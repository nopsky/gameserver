package main

import (
	"gameserver"
	//"log"
	"server/gate/service"
)

// func (this *TestHandle) ServiceHandle([]byte) ([]byte, error) {
// 	log.Println("这是业务逻辑处理")
// 	return nil, nil
// }

// func (this *TestHandle) SystemHandle([]byte) ([]byte, error) {
// 	log.Println("这是系统逻辑处理")
// 	return nil, nil
// }

// func (this *TestHandle) MultiCastHandle([]byte) ([]byte, error) {
// 	log.Println("这是广播逻辑处理")
// 	return nil, nil
// }

// func (this *TestHandle) ForwardHandle([]byte) ([]byte, error) {
// 	log.Println("这是转发逻辑处理")
// 	return nil, nil
// }

func main() {
	var ccs []gameserver.ClientConfig
	var sc gameserver.ServerConfig
	dh := NewHandle()
	sc.Name = "Gate Server"
	sc.Port = "8080"
	sc.MaxConn = 1024
	gs := gameserver.NewGameServer(ccs, sc, false, dh, dh)
	dh.Remote = gs.Cls
	dh.Local = service.Local
	gs.Start()
}
