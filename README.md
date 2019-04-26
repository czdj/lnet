# go游戏服务器网络库
## 层次结构
 每一层都定义了抽象的接口，只要实现相应的接口，可以方便的对每一层进行扩展。每一个接口都提供了一个基类实现，可以继承该基类，然后重写想要修改的方法即可。
- Transport层   
    负责网络相关的功能，目前支持TCP，WebSocket
    ```golang
        type ITransport interface {
            GetId() uint32
            GetLocalAddr() string
            GetRemoteAddr() string
            Listen() error
            Connect() error
            OnNewConnect(transport ITransport)
            Read()
            Write()
            Send(msgPkg IMessagePackage) error
            Close()
            OnClosed()
            IsStop() bool
            IsTimeout(tick *time.Timer) bool
        }
    ```
    
- Protocol层  
    负责协议解析，目前支持ProtoBuff，Gob编码协议
    ```golang
        type IProtocol interface {
        	Marshal(msg interface{}) ([]byte, error)
        	Unmarshal(data []byte, v interface{}) error
        }
    ```
- MsgHandle层   
    负责消息分发
    ```golang
        type IMsgHandle interface {
            RegisterMsg(tag uint32, msg interface{})
            NewMsg(tag uint32) interface{}
            GetMsgTag(msg interface{}) uint32
            SetProtocol(protocol IProtocol)
            GetProtocol() IProtocol
            CreateMessage(msgPkg IMessagePackage) interface{}
            CreateMessagePackage(msg interface{}) IMessagePackage
            Process(transport ITransport, msgPackage IMessagePackage)
            SetOnTransportClose(f func(transport ITransport))
            OnTransportClose(transport ITransport)
        }
    ```
## 样例
 构建一个服务只需组合上述三个层次的结构，即可定制需要的服务  
```golang
//服务器逻辑
func serverStart() {
    //选择协议
    protocol := &lprotocol.PbProtocol{}
    //选择消息处理函数
    msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
    //注册消息
    msgHandle.RegisterMsg(12, pb.GameItem{})
    msgHandle.RegisterMsg(13, pb.RpcReqData{})
    msgHandle.RegisterMsg(14, pb.RpcRspData{})
    //启动服务
    server := lserver.NewTcpServer("127.0.0.1:9000", msgHandle)
    server.Start()
}

//客户端逻辑
func clientStart() {
    //选择协议
    protocol := &lprotocol.PbProtocol{}
    //选择消息处理函数
    msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
    //注册消息
    msgHandle.RegisterMsg(12, pb.GameItem{})
    msgHandle.RegisterMsg(13, pb.RpcReqData{})
    msgHandle.RegisterMsg(14, pb.RpcRspData{})

    client := lclient.NewTcpClient("127.0.0.1:9000", msgHandle)
    client.Connect()

    msg := &pb.GameItem{Id: 1, Type: 2, Count: 3}
    for {
        client.Send(msg)
        time.Sleep(1 * time.Second)
    }
}
```
 
