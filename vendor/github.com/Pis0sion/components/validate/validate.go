package validate

import (
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
)

type validation struct {
	val   *validator.Validate
	data  interface{}
	trans ut.Translator
}

func NewValidator(data interface{}) *validation {

	val := validator.New()

	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	err := en.RegisterDefaultTranslations(val, trans)

	if err != nil {
		panic(err)
	}

	return &validation{
		val:   val,
		data:  data,
		trans: trans,
	}
}

func (v *validation) Validate() error {
	// validate policy
	err := v.val.Struct(v.data)
	if err == nil {
		return nil
	}

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if invalidValidationError, ok := err.(*validator.InvalidValidationError); ok {
		return invalidValidationError
	}

	// collect human-readable errors
	vErrors, _ := err.(validator.ValidationErrors)
	for _, vError := range vErrors {
		return vError
	}

	return err
}
