package request

type Login struct {
	Email        string `json:"email" validate:"required,email" field_error_info:"邮箱格式不正确"`
	Password     string `json:"password" validate:"required,pass" field_error_info:"密码格式应是8～32位且有大小写字母"`
	CheckCodeKey string `json:"checkCodeKey" validate:"required" field_error_info:"请通过正常方式访问"`
	CheckCode    string `json:"checkCode" validate:"required" field_error_info:"验证码不能为空"`
}
