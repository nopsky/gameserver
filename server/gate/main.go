package main

import (
	"gameserver"
	//"log"
	"server/gate/service"
)

func main() {
	cc := new(gameserver.ClientConfig)
	cc.GroupName = "hub"
	cc.ServerId = 1
	cc.ServerAddr = ":8081"
	var ccs []gameserver.ClientConfig
	ccs = append(ccs, *cc)
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
