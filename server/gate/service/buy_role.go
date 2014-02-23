package service

import (
	"log"
	"strconv"
)

import (
	"data"
	"message"
	"model"
)

type RoleBuyReq struct {
	Role_id uint8 //角色ID
}

type RoleBuyAck struct {
}

func init() {
	Local.Register("角色购买模块", message.MSG_BUYROLE, RoleBuy)
}

func RoleBuy(userInfo *model.UserInfo, reqData []byte) (ackData []byte, err error) {
	_req := &RoleBuyReq{}
	_ack := &RoleBuyAck{}
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
		reCode = message.ERR_BUYROLE
	}

	//说明还不拥有此角色
	if roleInfo.Id == 0 {
		err = userInfo.GetUserInfo(userInfo.Uid)
		//获取角色的信息（数值相关）
		cash := gamedata.GetInt("角色ID", strconv.Itoa(int(_req.Role_id)), "cash")
		diamond := gamedata.GetInt("角色ID", strconv.Itoa(int(_req.Role_id)), "diamond")
		log.Println("角色ID:", _req.Role_id, " cash:", cash)
		if userInfo.Cash > int32(cash) && userInfo.Diamond > int32(diamond) {
			err = userInfo.ChangeCash(userInfo.Uid, -cash)
			if err != nil {
				log.Println(userInfo.Uid, "购买角色信息出错", err)
				reCode = message.ERR_BUYROLE
			}
			err = userInfo.ChangeDiamond(userInfo.Uid, -diamond)
			if err != nil {
				log.Println(userInfo.Uid, "购买角色信息出错", err)
				reCode = message.ERR_BUYROLE
			}
			name := gamedata.GetString("角色ID", strconv.Itoa(int(_req.Role_id)), "角色名称")
			log.Println("角色ID:", _req.Role_id, " 角色名称:", name)
			err = roleInfo.AddRole(userInfo.Uid, _req.Role_id, name, 1, 1)
			if err != nil {
				log.Println(userInfo.Uid, "增加角色信息出错", err)
				reCode = message.ERR_BUYROLE
			}
		}
	}

	err = userInfo.ChangeRole(userInfo.Uid, _req.Role_id)
	if err != nil {
		log.Println(userInfo.Uid, "选择角色信息出错", err)
		reCode = message.ERR_BUYROLE
	}

	ackData = encode(userInfo.Uid, message.MSG_BUYROLE, reCode, _ack)
	return
}
