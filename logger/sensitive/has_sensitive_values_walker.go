package sensitive

import (
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

type hasSensitiveValuesWalker struct {
	hasSensitiveValues bool
	checkPackages      []string
}

var _ reflectwalk.StructWalker = &hasSensitiveValuesWalker{}

func (w *hasSensitiveValuesWalker) Struct(reflect.Value) error {
	return nil
}

func (w *hasSensitiveValuesWalker) StructField(f reflect.StructField, _ reflect.Value) error {
	if w.hasSensitiveValues {
		return reflectwalk.SkipEntry
	}

	if f.PkgPath != "" {
		return reflectwalk.SkipEntry
	}

	if f.Tag.Get(securedTag) == securedTagEnabledValue {
		w.hasSensitiveValues = true
		return reflectwalk.SkipEntry
	}

	if !mightHaveApplicationTypes(f.Type, w.checkPackages) {
		return reflectwalk.SkipEntry
	}

	return nil
}
