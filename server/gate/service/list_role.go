package service

import (
	"log"
)

import (
	"message"
	"model"
)

type RoleListReq struct {
}

type RoleListAck struct {
	List []model.RoleInfo
}

func init() {
	Local.Register("角色列表模块", message.MSG_LISTROLE, RoleList)
}

func RoleList(userInfo *model.UserInfo, reqData []byte) (ackData []byte, err error) {
	_req := &RoleListReq{}
	_ack := &RoleListAck{}
	err = decode(reqData, _req)
	if err != nil {
		log.Println("数据包格式不对", err)
		return nil, err
	}
	reCode := message.SUCCESS
	roleInfo := model.NewRoleInfo()
	_ack.List, err = roleInfo.GetRoleList(userInfo.Uid)

	if err != nil {
		log.Println("获取角色列表出错", err)
		return nil, err
	}
	ackData = encode(userInfo.Uid, message.MSG_LISTROLE, reCode, _ack)
	return
}
