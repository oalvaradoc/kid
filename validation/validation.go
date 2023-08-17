package validation

import (
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/validation/translations/th"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	zh_tw_translations "github.com/go-playground/validator/v10/translations/zh_tw"
	"reflect"
)

// CustomValidationRegisterFunc is a set configuration parameters of validation register function
type CustomValidationRegisterFunc struct {
	Tag                      string
	Func                     validator.Func
	CallValidationEvenIfNull bool
}

// RegisterDefaultTranslationsFn defines an interface that used to register the translation function
type RegisterDefaultTranslationsFn func(v *validator.Validate, trans ut.Translator) (err error)

// NewTranslator creates a new ut.Translator by ut.UniversalTranslator and locale
func NewTranslator(uni *ut.UniversalTranslator, locale string) ut.Translator {
	trans, _ := uni.GetTranslator(locale)
	return trans
}

// Tup2 represents a pair of validator and translator
type Tup2 struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

var validatorMap = map[string]Tup2{
	constant.LangZhCN: NewTup2(
		zh_translations.RegisterDefaultTranslations,
		NewTranslator(ut.New(zh.New()), "zh"),
	),
	constant.LangEnUS: NewTup2(
		en_translations.RegisterDefaultTranslations,
		NewTranslator(ut.New(en.New()), "en"),
	),
	constant.LangZhTW: NewTup2(
		zh_tw_translations.RegisterDefaultTranslations,
		NewTranslator(ut.New(en.New()), "zh"),
	),
	constant.LangThTH: NewTup2(
		th.RegisterDefaultTranslations,
		NewTranslator(ut.New(en.New()), "th"),
	),
}

// NewTup2 creates a tuple of validator and translator
func NewTup2(validatorFn func(v *validator.Validate, trans ut.Translator) (err error), translator ut.Translator) Tup2 {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})

	validatorFn(validate, translator)

	return Tup2{validate, translator}
}

// NewValidator creates a pair of validator and translator by user lang
func NewValidator(lang string) (*validator.Validate, ut.Translator) {
	tup2, ok := validatorMap[lang]
	if ok {
		return tup2.Validate, tup2.Trans
	}

	log.Warnsf("cannot found validator with language = %s using default language[%s] validator", lang, constant.LangEnUS)
	tup2 = validatorMap[constant.LangEnUS]

	return tup2.Validate, tup2.Trans
}
