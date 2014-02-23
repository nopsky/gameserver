package gameserver

import (
	"log"
	//"net"
)

//数据处理接口
type DataHandler interface {
	// ServiceHandle([]byte) ([]byte, error)
	// SystemHandle([]byte) ([]byte, error)
	// MultiCastHandle([]byte) ([]byte, error)
	// ForwardHandle([]byte) ([]byte, error)
	ClientHandle([]byte)
	ServerHandle([]byte, *Session) ([]byte, error)
}

type GameServer struct {
	Cls   *RemoteServices //远程处理服务
	ccs   []ClientConfig  //连接远程的配置参数
	sc    ServerConfig    //本地服务的配置参数
	cdh   DataHandler
	sdh   DataHandler
	Debug bool //是否是debug模式
}

type ClientConfig struct {
	GroupName  string
	ServerId   int32
	ServerAddr string
}

type ServerConfig struct {
	Name    string
	Port    string
	MaxConn int
}

func NewGameServer(ccs []ClientConfig, sc ServerConfig, debug bool, cdh DataHandler, sdh DataHandler) *GameServer {
	gs := new(GameServer)
	gs.ccs = ccs
	gs.sc = sc
	gs.Debug = debug
	gs.Cls = NewRemoteServices(len(ccs))
	gs.cdh = cdh
	gs.sdh = sdh
	return gs
}

//连接远程服务
func (gs GameServer) StartClient() {
	for _, cc := range gs.ccs {
		cli := NewClient(cc.GroupName, cc.ServerId, cc.ServerAddr, gs.cdh)
		err := cli.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		//注册client
		gs.Cls.Register(cc.GroupName, cc.ServerId, cli.Query)
	}
}

//启动本地服务
func (gs GameServer) StartServer() {
	srv := NewServer(gs.sc.Name, gs.sc.Port, gs.sc.MaxConn, gs.sdh)
	srv.Start()
}

//启动GameServer
func (gs GameServer) Start() {
	gs.StartClient()
	gs.StartServer()

}
