package main

import (
	"gameserver"
	//"log"
	"server/game/service"
)

func main() {
	var ccs []gameserver.ClientConfig
	cc := new(gameserver.ClientConfig)
	cc.GroupName = "hub"
	cc.ServerId = 1
	cc.ServerAddr = ":8081"
	ccs = append(ccs, *cc)
	var sc gameserver.ServerConfig
	dh := NewHandle()
	sc.Name = "Game Server"
	sc.Port = "8082"
	sc.MaxConn = 1024
	gs := gameserver.NewGameServer(ccs, sc, false, dh, dh)
	dh.Remote = gs.Cls
	dh.Local = service.Local
	gs.Start()
}
