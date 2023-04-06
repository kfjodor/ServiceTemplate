package sensitive

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

type hideSecuredFieldsWalker struct {
	depth         int
	skipDepth     int
	checkPackages []string
}

var _ reflectwalk.StructWalker = &hideSecuredFieldsWalker{} // ensure we implement appropriate interfaces
var _ reflectwalk.MapWalker = &hideSecuredFieldsWalker{}
var _ reflectwalk.SliceWalker = &hideSecuredFieldsWalker{}
var _ reflectwalk.ArrayWalker = &hideSecuredFieldsWalker{}
var _ reflectwalk.EnterExitWalker = &hideSecuredFieldsWalker{}

func (*hideSecuredFieldsWalker) Struct(reflect.Value) error {
	return nil
}

func (w *hideSecuredFieldsWalker) StructField(f reflect.StructField, v reflect.Value) error {
	if w.skipping() {
		return reflectwalk.SkipEntry
	}

	if f.PkgPath != "" {
		return reflectwalk.SkipEntry
	}

	if f.Tag.Get(securedTag) == securedTagEnabledValue {
		var err error

		switch i := v.Interface().(type) {
		case string:
			if i != "" && i != securedPlaceholder {
				err = safeSet(v, reflect.ValueOf(securedPlaceholder))
			}
		case *string:
			if i != nil && *i != "" && *i != securedPlaceholder {
				err = safeSet(v, reflect.ValueOf(strPtr(securedPlaceholder)))
			}
		default:
			zeroV := reflect.Zero(v.Type())
			if !reflect.DeepEqual(i, zeroV.Interface()) {
				err = safeSet(v, zeroV)
			}
		}

		if err != nil {
			return fmt.Errorf("field %+v is not modifiable - should never happen: %v", f, err)
		}

		return reflectwalk.SkipEntry // don't go inside secured structs when the *whole* struct marked secured
	}

	if isNil(v) {
		return reflectwalk.SkipEntry
	}

	if v.Type().Kind() == reflect.Interface || v.Type().Kind() == reflect.Ptr {
		securedV, err := valueOf(v.Interface(), w.checkPackages)
		if err != nil {
			return fmt.Errorf("failed to secure value of field %+v: %v", f, err)
		}

		if err := safeSet(v, reflect.ValueOf(securedV)); err != nil {
			// see comment above for another safeSet call
			return fmt.Errorf("field %+v is not modifiable - should never happen: %v", f, err)
		}
		return reflectwalk.SkipEntry
	}

	return nil
}

func (w *hideSecuredFieldsWalker) Map(m reflect.Value) error {
	if w.skipping() {
		return nil
	}
	return w.hideSecuredFieldsForMapValues(m)
}

func (w *hideSecuredFieldsWalker) MapElem(m, k, v reflect.Value) error {
	return nil // do nothing - necessary to implement reflect.MapWalker interface
}

func (w *hideSecuredFieldsWalker) Slice(s reflect.Value) error {
	if w.skipping() {
		return nil
	}
	return w.hideSecuredFieldsForSliceOrArrayValues(s)
}

func (w *hideSecuredFieldsWalker) SliceElem(i int, v reflect.Value) error {
	return nil // do nothing - necessary to implement reflect.SliceWalker interface
}

func (w *hideSecuredFieldsWalker) Array(s reflect.Value) error {
	if w.skipping() {
		return nil
	}
	return w.hideSecuredFieldsForSliceOrArrayValues(s)
}

func (w *hideSecuredFieldsWalker) ArrayElem(i int, v reflect.Value) error {
	return nil
}

func (w *hideSecuredFieldsWalker) hideSecuredFieldsForMapValues(m reflect.Value) error {
	keys := append([]reflect.Value{}, m.MapKeys()...)
	for _, key := range keys {
		value := m.MapIndex(key)

		if isNil(value) {
			continue
		}

		securedVal, err := valueOf(value.Interface(), w.checkPackages)
		if err != nil {
			return fmt.Errorf("failed to process map value at key %v: %v", key.Interface(), err)
		}

		m.SetMapIndex(key, reflect.ValueOf(securedVal))
	}

	w.skipCurrentEntryRecursive()
	return nil
}

func (w *hideSecuredFieldsWalker) hideSecuredFieldsForSliceOrArrayValues(s reflect.Value) error {
	for i := 0; i < s.Len(); i++ {
		value := s.Index(i)

		if isNil(value) {
			continue
		}

		securedVal, err := valueOf(value.Interface(), w.checkPackages)
		if err != nil {
			return fmt.Errorf("failed to process slice/array value of type %v at index %v: %v", value.Type(), i, err)
		}

		if err := safeSet(value, reflect.ValueOf(securedVal)); err != nil {
			return fmt.Errorf("failed to replace slice value at index %v: %v", i, err)
		}
	}

	w.skipCurrentEntryRecursive()
	return nil
}

func (w *hideSecuredFieldsWalker) Enter(_ reflectwalk.Location) error {
	w.depth++
	return nil
}

func (w *hideSecuredFieldsWalker) Exit(_ reflectwalk.Location) error {
	w.depth--

	if w.depth < w.skipDepth {
		w.skipDepth = 0
	}
	return nil
}

func (w *hideSecuredFieldsWalker) skipCurrentEntryRecursive() {
	w.skipDepth = w.depth
}

func (w *hideSecuredFieldsWalker) skipping() bool {
	return w.skipDepth > 0 && w.depth >= w.skipDepth
}

func safeSet(variable, value reflect.Value) error {
	if !variable.CanSet() {
		return fmt.Errorf("cannot set for variable of type %v", safeGetType(variable))
	}
	variable.Set(value)
	return nil
}

func safeGetType(v reflect.Value) reflect.Type {
	if !v.IsValid() {
		return nil
	}
	return v.Type()
}

func strPtr(s string) *string {
	return &s
}
