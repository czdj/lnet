package lnet

import (
	"fmt"
	"reflect"
)

var  DefMsgTimeout int  = 180

//要处理的消息，需要在此处注册
type MsgTypeMap struct {
	msgTypeMap map[uint16]reflect.Type
}

func (this *MsgTypeMap)Register(tag uint16,msg interface{}){
	msgType := reflect.TypeOf(msg)
	this.msgTypeMap[tag] = msgType
}

func (this *MsgTypeMap)NewMsg(tag uint16)interface{}{
	msgType,err := this.msgTypeMap[tag]
	if err == false{
		fmt.Println("Msg Type Err!")
		return nil
	}

	msg := reflect.New(msgType).Interface()

	return msg
}

var MsgTypeInfo MsgTypeMap = MsgTypeMap{msgTypeMap:make(map[uint16]reflect.Type)}


type NetType int
const (
	TCP NetType = iota
	UDP
	WebSocket
)

//网络包的格式为包头+包体组成，为TLV格式
type PakgeHead struct {
	Tag uint16
	Len uint16
}
type Pakge struct {
	head PakgeHead
	data []byte
}

//自定义消息类型
type MessageTest struct {
	Data string
}
