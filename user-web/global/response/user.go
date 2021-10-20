package response

type UserResponse struct {
	Id       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	Birthday uint64 `json:"birthday"`
}

type UserListResponse struct {
	BaseResponse
	Data []UserResponse `json:"data"`
}

func NewSuccessUserListResponse() UserListResponse {
	return NewUserListResponse(200, "OK")
}

func NewUserListResponse(code int, message string) UserListResponse {
	return UserListResponse{
		BaseResponse: BaseResponse{
			Code:    code,
			Message: message,
		},
		Data: make([]UserResponse, 0),
	}
}
