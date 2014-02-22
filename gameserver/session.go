package gameserver

import (
	"net"
	"time"
)

type Session struct {
	IP        net.IP
	SessionId uint64
	LoggedIn  bool
	KickOut   bool
	MQ        chan []byte

	// time related
	ConnectTime    time.Time
	LastPacketTime int64
	LastFlushTime  int64
	OpCount        int
}
