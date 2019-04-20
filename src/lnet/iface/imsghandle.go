package iface

//负责业务处理
type IMsgHandle interface {
	RegisterMsg(tag uint32, msg interface{})
	NewMsg(tag uint32) interface{}
	GetMsgTag(msg interface{}) uint32
	SetProtocol(protocol IProtocol)
	GetProtocol() IProtocol
	CreateMessage(msgPkg IMessagePackage) interface{}
	CreateMessagePackage(msg interface{}) IMessagePackage
	Process(transport ITransport, msgPackage IMessagePackage)
}
