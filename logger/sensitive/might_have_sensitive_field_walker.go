package sensitive

import (
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

type mightHaveSecuredFieldsWalker struct {
	mightHaveSecuredFields bool
	checkPackages          []string
}

var _ reflectwalk.StructWalker = &mightHaveSecuredFieldsWalker{} // ensure we implement appropriate interface

func (w *mightHaveSecuredFieldsWalker) Struct(reflect.Value) error {
	return nil
}

func (w *mightHaveSecuredFieldsWalker) StructField(f reflect.StructField, _ reflect.Value) error {
	if w.mightHaveSecuredFields {
		return reflectwalk.SkipEntry
	}

	if f.PkgPath != "" {
		return reflectwalk.SkipEntry
	}

	if f.Tag.Get(securedTag) == securedTagEnabledValue {
		w.mightHaveSecuredFields = true
		return reflectwalk.SkipEntry
	}

	if !mightHaveApplicationTypes(f.Type, w.checkPackages) {
		return reflectwalk.SkipEntry
	}

	kind := f.Type.Kind()

	if kind == reflect.Interface {
		w.mightHaveSecuredFields = true
		return nil
	}

	if kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map || kind == reflect.Ptr {
		valuesMightBeSecured, err := mightHaveSecuredFields(f.Type.Elem())
		if err != nil {
			return err
		}

		if valuesMightBeSecured {
			w.mightHaveSecuredFields = true
			return reflectwalk.SkipEntry
		}
	}

	return nil
}
