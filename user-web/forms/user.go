package forms

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password  string `form:"password" json:"password" binding:"required,min=8,max=20"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
	Answer    string `form:"answer" json:"answer" binging:"required,min=5,max=5"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}
