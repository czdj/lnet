package iface

type ITransportManager interface {
	Add(t ITransport)                  //添加连接
	Remove(t ITransport)               //删除连接
	Get(id uint32) (ITransport, error) //利用id获取连接
	Len() int32                        //获取当前连接
	ClearTransport()                   //删除并停止所有连接
}
