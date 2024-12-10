package defines

const (
	CAPTCHA          = "captcha:"
	CAPTCHA_TIMEOUT  = 5 * 60
	FIELD_ERROR_INFO = "field_error_info"
	PASSWORD_REGEX   = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d]{8,32}$`
)