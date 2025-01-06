package request

type AddFriend struct {
	FriendInfo string `json:"friendInfo" validate:"required" field_error_info:"好友信息不能为空"`
}
