package service

import (
	"errors"
	"fmt"
	"github.com/ugorji/go/codec"
	"lib/packet"
	"log"
	"model"
)

type funcType func(*model.UserInfo, []byte) ([]byte, error)

type LocalServices struct {
	servs map[int32]*serv
}

type serv struct {
	Name string
	Run  funcType
}

var Local = NewLocalService(10)

func NewLocalService(ServNum int) *LocalServices {
	s := new(LocalServices)
	s.servs = make(map[int32]*serv, ServNum)
	return s
}

func (this *LocalServices) Register(name string, msgid int32, run funcType) (err error) {
	if _serv, ok := this.servs[msgid]; !ok {
		s := new(serv)
		s.Name = name
		s.Run = run
		this.servs[msgid] = s
		log.Println("绑定服务:", name, " 对应的消息ID为:", msgid)
	} else {
		err = errors.New(fmt.Sprintf("此消息ID:%d 已经存在, name:%s\n", msgid, _serv.Name))
	}
	return
}

func (this *LocalServices) GetFunc(msgid int32) (run funcType, err error) {
	if _serv, ok := this.servs[msgid]; ok {
		run = _serv.Run
	} else {
		err = errors.New(fmt.Sprintf("此消息ID :%d 不存在\n", msgid))
	}
	return
}

func decode(reqData []byte, _req interface{}) (err error) {
	var mh codec.MsgpackHandle
	decode := codec.NewDecoderBytes(reqData, &mh)
	err = decode.Decode(_req)
	return
}

func encode(uid uint64, msgid int32, reCode int32, _ack interface{}) []byte {
	var mh codec.MsgpackHandle
	var out []byte
	mh.EncodeOptions.StructToArray = true
	encode := codec.NewEncoderBytes(&out, &mh)
	encode.Encode(_ack)
	writer := packet.Writer()
	//计算包的长度
	writer.WriteU64(uid)
	writer.WriteS32(msgid)
	writer.WriteS32(reCode)
	writer.WriteRawBytes(out)

	return writer.Data()
}
