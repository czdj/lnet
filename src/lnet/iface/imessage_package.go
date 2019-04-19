package iface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessagePackage interface {
	GetTag() uint32  //获取消息Tag
	GetLen() uint32  //获取消息数据段长度
	GetData() []byte //获取消息内容

	SetTag(uint32)  //设计消息ID
	SetData([]byte) //设计消息内容
	SetLen(uint32)  //设置消息数据段长度
}
