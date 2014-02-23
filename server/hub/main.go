package main

import (
	"gameserver"
	//"log"
	"server/hub/service"
)

func main() {
	var ccs []gameserver.ClientConfig
	var sc gameserver.ServerConfig
	dh := NewHandle()
	sc.Name = "Hub Server"
	sc.Port = "8081"
	sc.MaxConn = 1024
	gs := gameserver.NewGameServer(ccs, sc, false, dh, dh)
	dh.Remote = gs.Cls
	dh.Local = service.Local
	gs.Start()
}
