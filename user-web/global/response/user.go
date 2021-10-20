package response

type UserResponse struct {
	Id       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Birthday uint64 `json:"birthday"`
}

func NewSuccessUserListResponse(data interface{}) BaseResponse {
	return NewUserListResponse(200, "OK", data)
}

func NewFailedUserListResponse(code int, message string) BaseResponse {
	return NewUserListResponse(code, message, nil)
}

func NewUserListResponse(code int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
