package response

type UserResponse struct {
	Id       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Birthday uint64 `json:"birthday"`
}

type UserListResponse struct{
	Total uint32 `json:"total"`
	UserList []UserResponse `json:"user_list"`
}

func NewSuccessResponse(data interface{}) BaseResponse {
	return NewBaseResponse(200, "OK", data)
}

func NewFailedBaseResponse(code int, message string) BaseResponse {
	return NewBaseResponse(code, message, nil)
}

func NewBaseResponse(code int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
