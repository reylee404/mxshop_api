package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	mxValidator "mxshop_api/user-web/validator"
)

func MustInitValidators(tran *ut.Translator) {
	err := initValidators(tran)
	if err != nil {
		panic(err)
	}
}

func initValidators(tran *ut.Translator) error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("validator cast failed")
	}

	err := validate.RegisterValidation("mobile", mxValidator.ValidateMobile)
	if err != nil {
		return err
	}
	err = validate.RegisterTranslation("mobile", *tran,
		func(ut ut.Translator) error {
			return ut.Add("mobile", "手机号码不合法", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	if err != nil {
		return err
	}
	return nil
}
