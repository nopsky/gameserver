package msg

//请求的msgid

const SUCCESS int32 = 0
const (
	ERR_REGISTER int32 = 1 + iota
	ERR_LOGIN          = 2
)
