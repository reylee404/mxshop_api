package initialize

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"mxshop_api/user-web/global"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	mxValidator "mxshop_api/user-web/validator"
)

func MustInitValidators() {
	err := initValidators()
	if err != nil {
		panic(err)
	}
}

func initValidators() error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("validator cast failed")
	}

	err := validate.RegisterValidation("mobile", mxValidator.ValidateMobile)
	if err != nil {
		return err
	}
	err = validate.RegisterTranslation("mobile", global.Trans,
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
