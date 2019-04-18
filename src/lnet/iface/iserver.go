package iface

type IServer interface {
	Start()

	SetTransport(transport ITransport)
	//得到连接管理
	GetTransportMgr() ITransportManager
}
