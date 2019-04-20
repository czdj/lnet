package lmsghandle

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"reflect"
)

type BaseMsgHandle struct {
	msgTagTypeMap map[uint32]reflect.Type
	msgTypeTagMap map[reflect.Type]uint32
	protocol      iface.IProtocol
}

func NewBaseMsgHandle(protocol iface.IProtocol) *BaseMsgHandle {
	msgHandle := &BaseMsgHandle{msgTagTypeMap: make(map[uint32]reflect.Type), msgTypeTagMap: make(map[reflect.Type]uint32)}
	msgHandle.SetProtocol(protocol)

	return msgHandle
}
func (this *BaseMsgHandle) RegisterMsg(tag uint32, msg interface{}) {
	msgType := reflect.TypeOf(msg)
	this.msgTagTypeMap[tag] = msgType
	this.msgTypeTagMap[msgType] = tag
}

func (this *BaseMsgHandle) NewMsg(tag uint32) interface{} {
	msgType, err := this.msgTagTypeMap[tag]
	if err == false {
		lnet.Logger.Error("Msg Type Err!")
		return nil
	}

	msg := reflect.New(msgType).Interface()

	return msg
}

func (this *BaseMsgHandle) GetMsgTag(msg interface{}) uint32 {
	tag, err := this.msgTypeTagMap[reflect.TypeOf(msg).Elem()]
	if err == false {
		lnet.Logger.Error("Msg Type Err!")
		return 0
	}

	return tag
}

func (this *BaseMsgHandle) SetProtocol(protocol iface.IProtocol) {
	this.protocol = protocol
}

func (this *BaseMsgHandle) GetProtocol() iface.IProtocol {
	return this.protocol
}

func (this *BaseMsgHandle) CreateMessage(msgPkg iface.IMessagePackage) interface{} {
	msg := this.NewMsg(msgPkg.GetTag())
	err := this.protocol.Unmarshal(msgPkg.GetData(), msg)
	if err != nil {
		lnet.Logger.Error("msg Unmarshal err", zap.Any("err", err))
		return nil
	}

	return msg
}

func (this *BaseMsgHandle) CreateMessagePackage(msg interface{}) iface.IMessagePackage {
	encodeData, err := this.protocol.Marshal(msg)
	if err != nil {
		lnet.Logger.Error("msg marshal err", zap.Any("err", err))
		return nil
	}
	tag := this.GetMsgTag(msg)
	msgPkg := lnet.NewMsgPackage(tag, encodeData)

	return msgPkg
}

func (this *BaseMsgHandle) Process(itransport iface.ITransport, msgPackage iface.IMessagePackage) {
	msg := this.CreateMessage(msgPackage)

	lnet.Logger.Info("process msg", zap.Any("RemoteAddr", itransport.GetRemoteAddr()), zap.Any("msg", msg))

	msgPkg := this.CreateMessagePackage(msg)

	itransport.Send(msgPkg)
}
