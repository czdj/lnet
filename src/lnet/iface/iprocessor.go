package iface

//负责业务处理
type IProcessor  interface {
	Process(transport ITransport, msg interface{})
}

