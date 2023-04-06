package sensitive

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/reflectwalk"
	"github.com/pkg/errors"
)

var (
	mightHaveSecuredFieldsCache = map[reflect.Type]bool{}
	cacheMutex                  sync.RWMutex
)

type Value struct {
	v             interface{}
	checkPackages []string
}

const (
	securedTag             = "hide"
	securedTagEnabledValue = "true"
	securedPlaceholder     = "<hidden>"
)

func NewValue(v interface{}, checkPackages []string) Value {
	return Value{v: v, checkPackages: checkPackages}
}

func (v Value) MarshalJSON() ([]byte, error) {
	safeV, err := valueOf(v.v, v.checkPackages)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to obtain safe value")
	}
	return json.Marshal(safeV)
}

func valueOf(in interface{}, checkPackages []string) (safe interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			safe = nil
			err = fmt.Errorf("panic while securing value of type %T: %v", in, e)
		}
	}()

	if isNil(reflect.ValueOf(in)) {
		return in, nil
	}

	nonPtrValueType := reflect.Indirect(reflect.ValueOf(in)).Type()

	if !mightHaveApplicationTypes(nonPtrValueType, checkPackages) {
		return in, nil
	}

	if reflect.ValueOf(in).Kind() == reflect.Ptr && nonPtrValueType.Kind() == reflect.Interface {
		vElem := reflect.ValueOf(in).Elem()
		if isNil(vElem) {
			return in, nil
		}

		safeV, verr := valueOf(vElem.Interface(), checkPackages)
		if verr != nil {
			return nil, verr
		}

		ptr := reflect.New(nonPtrValueType)
		ptr.Elem().Set(reflect.ValueOf(safeV))

		return ptr.Interface(), nil
	}

	mightHaveSecuredFields, err := mightHaveSecuredFields(nonPtrValueType)
	if err != nil {
		return nil, fmt.Errorf("failed to check for secured fields of type %T", in)
	}

	if !mightHaveSecuredFields {
		return in, nil
	}

	hasValues, err := hasValues(in, checkPackages)
	if err != nil {
		return nil, fmt.Errorf("failed to check for secured fields of value %#v", in)
	}

	if !hasValues {
		return in, nil
	}

	moveToPointer := reflect.ValueOf(in).Kind() == reflect.Struct || reflect.ValueOf(in).Kind() == reflect.Array

	if moveToPointer {
		ptr := reflect.New(nonPtrValueType)
		ptr.Elem().Set(reflect.ValueOf(in))
		in = ptr.Interface()
	}

	vCopy, err := copystructure.Copy(in)
	if err != nil {
		return nil, fmt.Errorf("failed to deep copy value of type %T", in)
	}

	if err := hideValues(vCopy, checkPackages); err != nil {
		return nil, fmt.Errorf("failed to hide secured values from value of type %T: %v", vCopy, err)
	}

	if !moveToPointer {
		return vCopy, nil
	}

	return reflect.Indirect(reflect.ValueOf(vCopy)).Interface(), nil
}

func mightHaveSecuredFields(t reflect.Type) (bool, error) {
	cacheMutex.RLock()
	mightHaveSecuredFields, ok := mightHaveSecuredFieldsCache[t]
	cacheMutex.RUnlock()

	if ok {
		return mightHaveSecuredFields, nil
	}

	cacheMutex.Lock()
	mightHaveSecuredFieldsCache[t] = false
	cacheMutex.Unlock()

	mightHaveSecuredFields, err := calculateMightHaveSecuredFields(t)
	if err != nil {
		return false, err
	}

	cacheMutex.Lock()
	mightHaveSecuredFieldsCache[t] = mightHaveSecuredFields
	cacheMutex.Unlock()

	return mightHaveSecuredFields, nil
}

func calculateMightHaveSecuredFields(t reflect.Type) (bool, error) {
	kind := t.Kind()

	if kind == reflect.Interface {
		return true, nil // anything could be there on runtime - so it _might_ have secured tags
	}

	if kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map || kind == reflect.Ptr {
		return mightHaveSecuredFields(t.Elem())
	}

	walker := mightHaveSecuredFieldsWalker{}

	if err := reflectwalk.Walk(reflect.Zero(t).Interface(), &walker); err != nil {
		return false, err
	}

	return walker.mightHaveSecuredFields, nil
}

func hasValues(in interface{}, checkPackages []string) (bool, error) {
	walker := hasSensitiveValuesWalker{checkPackages: checkPackages}

	if err := reflectwalk.Walk(in, &walker); err != nil {
		return false, err
	}

	return walker.hasSensitiveValues, nil
}

func hideValues(in interface{}, checkPackages []string) error {
	return reflectwalk.Walk(in, &hideSecuredFieldsWalker{checkPackages: checkPackages})
}

func isNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface || v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		return v.IsNil()
	}
	return false
}

func mightHaveApplicationTypes(t reflect.Type, checkPackages []string) bool {
	pkgPath := t.PkgPath()

	if pkgPath == "" || checkPackages == nil {
		return true
	}

	for _, pkg := range checkPackages {
		if pkgPath == pkg {
			return true
		}
	}

	return false
}
