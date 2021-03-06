package lnet

const (
	MAX_CONN    = 10000
	MAX_PKG_LEN = 1024 * 1024
)

var DefMsgTimeout int = 180 //秒

////要处理的消息，需要在此处注册
//type MsgTypeMap struct {
//	msgTagTypeMap map[uint32]reflect.Type
//	msgTypeTagMap map[reflect.Type]uint32
//}
//
//func (this *MsgTypeMap) Register(tag uint32, msg interface{}) {
//	msgType := reflect.TypeOf(msg)
//	this.msgTagTypeMap[tag] = msgType
//	this.msgTypeTagMap[msgType] = tag
//}
//
//func (this *MsgTypeMap) NewMsg(tag uint32) interface{} {
//	msgType, err := this.msgTagTypeMap[tag]
//	if err == false {
//		fmt.Println("Msg Type Err!")
//		return nil
//	}
//
//	msg := reflect.New(msgType).Interface()
//
//	return msg
//}
//
//func (this *MsgTypeMap) Tag(msg interface{}) uint32 {
//	tag, err := this.msgTypeTagMap[reflect.TypeOf(msg).Elem()]
//	if err == false {
//		fmt.Println("Msg Type Err!")
//		return 0
//	}
//
//	return tag
//}
//
//var MsgTypeInfo MsgTypeMap = MsgTypeMap{msgTagTypeMap: make(map[uint32]reflect.Type), msgTypeTagMap: make(map[reflect.Type]uint32)}

type NetType int

const (
	TCP NetType = iota
	UDP
	WebSocket
)
