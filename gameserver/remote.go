package gameserver

import (
	"errors"
	"fmt"
)

type funcType func([]byte) ([]byte, error)

type RemoteServices struct {
	servs map[string]map[int32]*rserv // 服务标识->id->服务接口, 支持一个服务下面多台服务器集群
}

type rserv struct {
	Run funcType
}

func NewRemoteServices(ServNum int) *RemoteServices {
	s := new(RemoteServices)
	s.servs = make(map[string]map[int32]*rserv, ServNum)
	return s
}

func (this *RemoteServices) Register(groupName string, serverId int32, run funcType) (err error) {
	if _, ok := this.servs[groupName][serverId]; !ok {
		s := new(rserv)
		s.Run = run
		this.servs[groupName][serverId] = s
	} else {
		err = errors.New(fmt.Sprintf("此服务ID: %d 已经存在, name: %s \n", serverId, groupName))
	}
	return
}

func (this *RemoteServices) GetFunc(groupName string, serverId int32) (run funcType, err error) {
	if _serv, ok := this.servs[groupName][serverId]; ok {
		run = _serv.Run
	} else {
		err = errors.New(fmt.Sprintf("当前groupName : %s serverId: %d 不存在\n", groupName, serverId))
	}
	return
}
