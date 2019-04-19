package iface

//负责解析协议
type IProtocol interface {
	Marshal(msg interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}
