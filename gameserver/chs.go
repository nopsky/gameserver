package gameserver

import (
	"errors"
	"fmt"
	"sync"
)

var Chs = make(map[uint64]chan []byte, 1024)

var chs_lock sync.RWMutex

//注册Conn
func AddConn(id uint64, ch chan []byte) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	Chs[id] = ch
}

//移除某个Conn
func RemoveConn(id uint64) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	delete(Chs, id)
}

//获取某个Conn
func GetConn(id uint64) (ch chan []byte, err error) {
	chs_lock.RLock()
	defer chs_lock.RUnlock()
	ch, ok := Chs[id]
	if !ok {
		err = errors.New(fmt.Sprintf("uid %x 没有对应的 Session", id))
		return
	}
	return
}
