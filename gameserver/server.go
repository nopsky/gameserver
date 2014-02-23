package gameserver

import (
	"encoding/binary"
	//"errors"
	"hash/crc32"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	//"sync/atomic"
)

// const (
// 	SRC_CLIENT = 0 //客户端请求
// 	SRC_IPC    = 1 //内网请求
// )

type Server struct {
	Name    string //服务名称
	Port    string //服务端口
	Conn    *Connect
	MaxConn int         //最大连接数
	DH      DataHandler //数据处理接口
}

func NewServer(name string, port string, max int, dh DataHandler) *Server {
	s := new(Server)
	s.Name = name
	s.Port = port
	s.MaxConn = max
	s.DH = dh
	return s
}

func (s Server) Start() {
	log.Println("启动服务:", s.Name)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+s.Port)

	if err != nil {
		log.Println(err)
		return
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		c := NewConnect(conn, s.DH)
		//封禁IP操作
		//IP := net.ParseIP(conn.RemoteAddr().String())
		//if !IsBanned(IP) {
		go c.StartAgent()
		//} else {
		//conn.Close()
		//}
	}
}

type Connect struct {
	conn net.Conn
	wg   sync.WaitGroup
	exit chan bool   //退出信号
	send chan []byte //发送数据通道
	recv chan []byte //接受数据通道
	dh   DataHandler //数据处理协议函数
	sess *Session    //当前会话的Seesion信息
}

func NewConnect(conn net.Conn, dh DataHandler) *Connect {
	return &Connect{
		conn: conn,
		exit: make(chan bool),
		dh:   dh,
		send: make(chan []byte, 1),
		sess: new(Session),
	}
}

//启动处理协程,没一个客户端对应一个recv,send
func (c *Connect) StartAgent() {
	c.wg.Add(2)
	go func() {
		defer c.wg.Done()
		c._recv()
	}()
	go func() {
		defer c.wg.Done()
		c._send()
	}()
}

//客户端关闭时,退出recv,send
func (c *Connect) Close() {
	close(c.exit)
	c.conn.Close()
	c.wg.Wait()
}

func (c *Connect) _recv() (err error) {
	defer func() {
		if err != nil {
			select {
			case <-c.exit:
				err = nil
			default:

			}
		}
		//记录当前退出的Session日志
	}()
	c.sess.IP = net.ParseIP(strings.Split(c.conn.RemoteAddr().String(), ":")[0])
	c.sess.ConnectTime = time.Now()
	c.sess.LastPacketTime = time.Now().Unix()
	c.sess.KickOut = false
	c.sess.MQ = make(chan []byte, 1024)
	//1秒的时钟
	timer := time.NewTicker(1 * time.Second)
	for {
		var data []byte
		select {
		case <-c.exit: //接受到退出信息时,goroutine结束
			return nil
		case <-timer.C: //时钟,用来判断连接是否存活
		case data = <-c.sess.MQ: //接受向此连接发送的数据
			ackData, err := c.dh.ServerHandle(data, c.sess)

			if err != nil {
				log.Println(err)
				continue
			}

			if ackData != nil {
				c.send <- PacketData(0, ackData)
			}
		}

		header := make([]byte, 14)

		n, err := io.ReadFull(c.conn, header)

		if n == 0 && err == io.EOF {
			log.Println("客户端断开连接")
			break
		} else if err != nil {
			log.Println("接受数据出错:", err)
		}
		//数据包长度
		size := binary.BigEndian.Uint16(header)

		//crc值
		crc1 := binary.BigEndian.Uint32(header)

		data = make([]byte, size)

		n, err = io.ReadFull(c.conn, data)

		if uint16(n) != size {
			log.Println("数据包长度不正确", n, "!=", size)
			continue
		}

		if err != nil {
			log.Println("读取数据出错:", err)
			continue
		}

		crc2 := crc32.Checksum(data, crc32.IEEETable)

		if crc1 != crc2 {
			log.Println("crc 数据验证不正确: ", crc1, " != ", crc2)
			continue
		}

		seqId := binary.BigEndian.Uint64(header[6:])

		ackData, err := c.dh.ServerHandle(data, c.sess)

		if err != nil {
			log.Println(err)
			continue
		}

		if ackData != nil {
			c.send <- PacketData(seqId, ackData)
		}
		//指定数据类型
		// reader := packet.reader(data)

		// ptype, _ := reader.ReadU16()

		// body, err := reader.ReadAtLeast()

		// switch ptype {
		// case PTYPE_SERVICE:
		// 	//业务逻辑处理
		// 	ackData, err := c.dh.ServiceHandle(body)
		// case PTYPE_SYSTEM:
		// 	//系统逻辑处理
		// 	ackData, err := c.dh.SystemHandle(body)
		// case PTYPE_MULTICAST:
		// 	//广播逻辑处理
		// 	ackData, err := c.dh.MultiCastHandle(body)
		// case PTYPE_FORWARD:
		// 	//转发逻辑处理
		// 	ackData, err := c.dh.ForwardHandle(body)
		// default:
		// 	log.Println("错误的PTYPE类型:", ptype)
		// }
	}
	//关闭连接
	log.Println("退出此连接")
	c.Close()
	return nil
}

//发送数据
func (c *Connect) _send() {
	for {
		select {
		case <-c.exit:
			return
		case data := <-c.send:
			if _, err := c.conn.Write(data); err != nil {
				log.Println("发送数据出错", data)
				continue
			}
		}
	}
}
