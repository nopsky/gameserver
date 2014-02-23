package service

import (
	"log"
)

import (
	//"data"
	"message"
	"model"
)

type RoleLeveUpReq struct {
	Role_id uint8 //角色ID
}

type RoleLeveUpReqAck struct {
}

func init() {
	Local.Register("角色升级模块", message.MSG_LEVELUPROLE, RoleLevelUp)
}

func RoleLevelUp(userInfo *model.UserInfo, reqData []byte) (ackData []byte, err error) {
	_req := &RoleLeveUpReq{}
	_ack := &RoleLeveUpReqAck{}
	err = decode(reqData, _req)
	if err != nil {
		log.Println("数据包格式不对", err)
		return nil, err
	}
	reCode := message.SUCCESS

	ackData = encode(userInfo.Uid, message.MSG_LEVELUPROLE, reCode, _ack)
	return
}
