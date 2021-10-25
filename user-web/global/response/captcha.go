package response

type CaptchaResponse struct {
	CaptchaId   string `json:"captcha_id"`
	CaptchaPath string `json:"captcha_path"`
}
