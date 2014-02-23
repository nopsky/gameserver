package gameserver

import (
	"model"
	"net"
	"time"
)

type Session struct {
	IP      net.IP
	User    *model.UserInfo
	KickOut bool
	MQ      chan []byte

	// time related
	ConnectTime    time.Time
	LastPacketTime int64
	LastFlushTime  int64
	OpCount        int
}
