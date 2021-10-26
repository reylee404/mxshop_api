package initialize

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

func MustInitTrans(local string) ut.Translator {
	tran, err := InitTrans(local)
	if err != nil {
		panic(err)
	}
	return tran
}

func InitTrans(local string) (ut.Translator, error) {
	tran, err := initTrans(local)
	if err != nil {
		return nil, err
	}
	return tran, nil
}

func initTrans(local string) (ut.Translator, error) {
	var t ut.Translator
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
		t, ok = uni.GetTranslator(local)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}

		var err error
		switch local {
		case "en":
			err = entranslations.RegisterDefaultTranslations(v, t)
		case "zh":
			err = zhtranslations.RegisterDefaultTranslations(v, t)
		default:
			err = entranslations.RegisterDefaultTranslations(v, t)
		}
		return t, err
	}
	return nil, errors.New("validate cast failed")
}
