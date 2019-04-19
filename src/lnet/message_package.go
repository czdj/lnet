package lnet

//TLV格式
type MessagePackage struct {
	Tag  uint32 //消息的Tag
	Len  uint32 //消息的长度
	Data []byte //消息的内容
}

//创建一个Message消息包
func NewMsgPackage(tag uint32, data []byte) *MessagePackage {
	return &MessagePackage{
		Tag:  tag,
		Len:  uint32(len(data)),
		Data: data,
	}
}

//获取消息数据段长度
func (msg *MessagePackage) GetLen() uint32 {
	return msg.Len
}

//获取消息Tag
func (msg *MessagePackage) GetTag() uint32 {
	return msg.Tag
}

//获取消息内容
func (msg *MessagePackage) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *MessagePackage) SetLen(len uint32) {
	msg.Len = len
}

//设置消息ID
func (msg *MessagePackage) SetTag(tag uint32) {
	msg.Tag = tag
}

//设置消息内容
func (msg *MessagePackage) SetData(data []byte) {
	msg.Data = data
}
