package service

import (
	"log"
)

import (
	"message"
	"model"
)

type LoginReq struct {
	Username string
	Password string
}

type LoginAck struct {
	//用户信息
	UserInfo *model.UserInfo
	//角色信息
	RoleInfo *model.RoleInfo
}

func init() {
	Local.Register("用户登录模块", message.MSG_LOGIN, Login)
}

func Login(userInfo *model.UserInfo, reqData []byte) (ackData []byte, err error) {
	log.Println("登录数据包是:", reqData)
	_req := &LoginReq{}
	_ack := &LoginAck{}
	err = decode(reqData, _req)
	if err != nil {
		log.Println("数据包格式不对", err)
		return nil, err
	}
	reCode := message.SUCCESS
	roleInfo := model.NewRoleInfo()

	err = userInfo.CheckLogin(_req.Username, _req.Password)
	if err != nil {
		reCode = message.ERR_LOGIN
		err = nil
	} else {
		//根据用户uid得到角色信息
		err = roleInfo.GetRoleInfo(userInfo.Uid, userInfo.Role_Id)
		log.Println(err, roleInfo)
		if err != nil {
			log.Println("获取角色信息出错")
			return nil, err
		}
	}

	//建立当前会话信息

	log.Println("登陆返回的uid为:", userInfo.Uid)
	_ack.UserInfo = userInfo
	_ack.RoleInfo = roleInfo
	ackData = encode(userInfo.Uid, message.MSG_LOGIN, reCode, _ack)
	return
}
