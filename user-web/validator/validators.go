package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var mobilRegexp = regexp.MustCompile(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`)

func ValidateMobile(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return mobilRegexp.MatchString(s)
}
