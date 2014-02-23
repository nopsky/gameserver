package service

import (
	"log"
)

import (
	"message"
	"model"
)

type RoleReq struct {
	Role_id uint8 //角色ID
}

type RoleAck struct {
}

func init() {
	Local.Register("选择角色模块", message.MSG_SELECTROLE, SelectRole)
}

func SelectRole(userInfo *model.UserInfo, reqData []byte) (ackData []byte, err error) {
	_req := &RoleReq{}
	_ack := &RoleAck{}
	err = decode(reqData, _req)
	if err != nil {
		log.Println("数据包格式不对", err)
		return nil, err
	}
	reCode := message.SUCCESS
	roleInfo := model.NewRoleInfo()
	err = roleInfo.GetRoleInfo(userInfo.Uid, _req.Role_id)
	if err != nil {
		log.Println(userInfo.Uid, "查询角色信息出错", err)
		return nil, err
	}

	if roleInfo.Id == 0 {
		reCode = message.ERR_SELECTROLE
	}
	err = userInfo.ChangeRole(userInfo.Uid, _req.Role_id)
	if err != nil {
		log.Println(userInfo.Uid, "选择角色信息出错", err)
		return nil, err
	}
	ackData = encode(userInfo.Uid, message.MSG_SELECTROLE, reCode, _ack)
	return
}
