package service

import (
	"log"
)

import (
	gs "gameserver"
	"message"
	"model"
)

type RegisterReq struct {
	Username string
	Password string
}

type RegisterAck struct {
	//用户信息
	UserInfo *model.UserInfo
	//角色信息
	RoleInfo *model.RoleInfo
}

func init() {
	Local.Register("用户注册模块", msg.MSG_REGISTER, Reg)
}

func Reg(sess *gs.Session, reqData []byte) (ackData []byte, err error) {
	_req := &RegisterReq{}
	_ack := &RegisterAck{}
	err = decode(reqData, _req)
	if err != nil {
		log.Println("数据包格式不对", err)
		return nil, err
	}
	reCode := msg.SUCCESS
	userInfo := model.NewUserInfo()
	roleInfo := model.NewRoleInfo()
	err = userInfo.AddUser(_req.Username, _req.Password)
	if err != nil {
		reCode = msg.ERR_REGISTER
		err = nil
	} else {
		//根据用户uid得到角色信息
		err := roleInfo.AddRole(userInfo.Uid, 1, "summer", 1, 1)
		if err != nil {
			log.Println("增加默认角色出错", err)
			return nil, err
		}
		sess.SessionId = userInfo.Uid
		sess.LoggedIn = true
	}

	log.Println("登陆返回的uid为:", userInfo.Uid)
	_ack.UserInfo = userInfo
	_ack.RoleInfo = roleInfo
	ackData = encode(userInfo.Uid, msg.MSG_REGISTER, reCode, _ack)
	return
}
