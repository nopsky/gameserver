package main

import (
	"errors"
	"fmt"
	//"log"
)

import (
	gs "gameserver"
	"lib/packet"
	"message"
	"server/gate/service"
)

type Handle struct {
	Local  *service.LocalServices
	Remote *gs.RemoteServices
}

func NewHandle() *Handle {
	return new(Handle)
}

func (this *Handle) ServerHandle(data []byte, sess *gs.Session) (ackData []byte, err error) {
	reader := packet.Reader(data)

	//读取用户uid
	uid, err := reader.ReadU64()

	if err != nil {
		errstr := fmt.Sprintf("读取用户UID出错")
		err = errors.New(errstr)
		return
	}

	if uid != sess.SessionId {
		errstr := fmt.Sprintf("用户UID不正确,非法请求, uid:%d != sess.Uid:%d", uid, sess.SessionId)
		err = errors.New(errstr)
		return
	}

	//读取消息ID
	msgid, err := reader.ReadS32()

	if err != nil {
		errstr := fmt.Sprintf("读取消息ID出错")
		err = errors.New(errstr)
		return
	}

	//读取MsgPack的数据
	reqData, err := reader.ReadAtLeast()

	if err != nil {
		errstr := fmt.Sprintf("读取数据包内容出错")
		err = errors.New(errstr)
		return
	}

	_handle, err := this.Local.GetFunc(msgid)

	if err == nil {
		ackData, err = _handle(sess, reqData)
		if err != nil {
			return
		}
		if msgid == msg.MSG_REGISTER || msgid == msg.MSG_LOGIN {
			gs.AddConn(sess.SessionId, sess.MQ)
		} else if msgid == msg.MSG_LOGOUT {
			gs.RemoveConn(uid)
		}
		return
	}
	groupName, serverId, err := this.remote_hash(uid, msgid)
	if err != nil {
		return
	}
	_rhandle, err := this.Remote.GetFunc(groupName, serverId)

	if err == nil {
		//转发给其他服务器处理
		ackData, err = _rhandle(data)
		return
	}
	return
}

func (this *Handle) ClientHandle(data []byte) {

	// log.Println("处理客户端clientHandle")
	// reader := packet.Reader(data)

	// //读取用户uid
	// uid, err := reader.ReadU64()

	// if err != nil {
	// 	errstr := fmt.Sprintf("读取用户UID出错")
	// 	err = errors.New(errstr)
	// 	return
	// }

	// msgid := reader.ReadS32()

	// _handle, err := this.Local.GetFunc(msgid)
	// if _handle != nil {
	// 	_, err := _handle(reqData)
	// } else {
	// 	//转发给server
	// 	ch := gs.GetConn(uid)
	// 	ch <- data
	// }
}

func (this Handle) remote_hash(uid uint64, msgid int32) (groupName string, serverId int32, err error) {
	if uid == 0 || msgid == 0 {
		errstr := fmt.Sprintf("读取返回值的UID出错")
		err = errors.New(errstr)
		return
	}
	//先直接返回
	return "hub", 1, nil

}
