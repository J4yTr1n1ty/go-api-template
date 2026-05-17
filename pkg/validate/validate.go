package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// FieldErrors maps JSON field names → human-readable messages.
// Keys use dot-separated JSON tag names for nested structs.
type FieldErrors map[string]string

func (fe FieldErrors) Error() string { return "validation failed" }

var (
	v     *validator.Validate
	trans ut.Translator
)

func init() {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	v = validator.New(validator.WithRequiredStructEnabled())

	// Use JSON tag names in field error keys so callers see e.g. "first_name" not "FirstName".
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	en_translations.RegisterDefaultTranslations(v, trans) //nolint:errcheck
}

// Struct validates s and returns FieldErrors (keyed by JSON tag name) or nil.
func Struct(s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	fe := make(FieldErrors, len(errs))
	for _, e := range errs {
		// Namespace: "StructName.field.nested" → strip first segment → "field.nested"
		ns := e.Namespace()
		_, key, found := strings.Cut(ns, ".")
		if !found {
			key = ns
		}
		fe[key] = e.Translate(trans)
	}
	return fe
}
