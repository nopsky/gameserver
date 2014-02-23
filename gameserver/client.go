package gameserver

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"log"
	"net"
	//"os"
	"sync"
	"sync/atomic"
)

// type Client struct {
// 	Clients map[string][int]*Connect
// }

// func NewClient() *Client {
// 	client := new(Client)
// 	client.Clients = make(map[string][int]*Connect, 5)
// 	return client
// }

// //新连接
// func (cl *Client) Connect(name string, serverId int, serverAddr string) {

// }

// //关闭连接
// func (cl *Client) CloseConnect(name string, serverId string) {

// }

var ERR_EXIT = errors.New("连接退出")

type Client struct {
	Name    string
	SerId   int32
	SerAddr string
	//Type       int //当前客户端的类型: 0:玩家客户端 1:内网客户端
	conn       net.Conn
	wg         sync.WaitGroup
	mutex      sync.Mutex
	sess       *Session
	exit       chan bool   //退出信号
	send       chan []byte //发送数据通道
	dh         DataHandler
	seq        uint64 //会话ID
	applicants map[uint64]chan []byte
}

func NewClient(name string, serId int32, serAddr string, dh DataHandler) *Client {
	client := new(Client)
	client.Name = name
	client.SerId = serId
	client.SerAddr = serAddr
	//client.Type = SRC_IPC
	client.sess = new(Session)
	client.exit = make(chan bool)
	client.send = make(chan []byte, 1024)
	client.dh = dh
	client.applicants = make(map[uint64]chan []byte, 1024)
	return client
}

func (c *Client) Connect() (err error) {
	log.Println("Connecting to", c.Name)

	addr, err := net.ResolveTCPAddr("tcp", c.SerAddr)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
	c.conn = conn

	go c.StartAgent()
	return
}

// type Client struct {
// 	conn       net.Conn
// 	wg         sync.WaitGroup
// 	mutex      sync.Mutex
// 	applicants map[uint64]chan []byte
// 	exit       chan bool             //退出信号
// 	send       chan []byte           //发送数据通道
// 	seq        uint64                //会话ID
// 	maxch      int                   //每个客户端最大的并发通道
// 	dh         GameServer.DataHandle //数据处理协议函数
// 	eh         GameServer.ErrHandle  //错误处理协议函数
// 	sess       Session
// 	//recv       chan chan []byte      //接受数据通道的通道(每一个会话对应一个接受通道)
// }

//启动处理协程,每一个客户端对应一个recv,send
func (c *Client) StartAgent() {
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
func (c *Client) Close() {
	close(c.exit)
	c.conn.Close()
	c.wg.Wait()
}

func (c *Client) _recv() (err error) {
	defer func() {
		if err != nil {
			select {
			case <-c.exit:
				err = nil
			default:

			}
		}
	}()
	for {
		select {
		case <-c.exit:
			//接受到退出信息时,goroutine结束
			return nil
		default:
			break
		}

		var data []byte

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
		crc1 := binary.BigEndian.Uint32(header[2:6])

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

		if seqId > 0 {
			ch, ok := c.popApplicant(seqId)
			if ok {
				ch <- data
				continue
			}
		}
		//处理数据或者转发给server的接受通道
		c.dh.ClientHandle(data)
	}
	return nil
}

//发送数据
func (c *Client) _send() {
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

//需要响应的请求
func (c *Client) Query(data []byte) (res []byte, err error) {

	//指定接受响应数据的通道
	var ch chan []byte

	seq := c.newSeqId()

	//发送数据
	if err := c.Send(data, seq); err != nil {
		return nil, err
	}

	c.addApplicant(seq, ch)
	for {
		select {
		case <-c.exit:
			return nil, ERR_EXIT
		case res = <-ch:
			break
		}
	}
	return res, nil
}

//不需要响应的请求
func (c *Client) Send(data []byte, seq uint64) (err error) {
	select {
	case <-c.exit:
		return ERR_EXIT
	case c.send <- data:
		break
	}
	return nil
}

//得到新的会话ID
func (c *Client) newSeqId() uint64 {
	return atomic.AddUint64(&c.seq, 1)
}

//绑定通道和会话ID
func (c *Client) addApplicant(seq uint64, ch chan []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.applicants[seq] = ch
}

//删除通道和会话ID的绑定关系
func (c *Client) popApplicant(seq uint64) (chan []byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ch, ok := c.applicants[seq]
	if !ok {
		return nil, false
	}
	delete(c.applicants, seq)
	return ch, true
}
