package message

//请求的msgid
const (
	MSG_REGISTER int32 = 10001 + iota
	MSG_LOGIN
	MSG_LOGOUT
	MSG_SELECTROLE
	MSG_BUYROLE
	MSG_LEVELUPROLE
	MSG_LISTROLE
	//掉线
	MDropping
)
