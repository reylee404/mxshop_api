package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"mxshop_api/user-web/global"
	"reflect"
	"strings"
)

func MustInitTrans(local string) {
	err := InitTrans(local)
	if err != nil {
		panic(err)
	}
}

func InitTrans(local string) error {
	err := initTrans(local)
	if err != nil {
		return err
	}
	return nil
}

func initTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}

		var err error
		switch local {
		case "en":
			err = entranslations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			err = zhtranslations.RegisterDefaultTranslations(v, global.Trans)
		default:
			err = entranslations.RegisterDefaultTranslations(v, global.Trans)
		}
		return err
	}
	return
}
