package lnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"lnet/iface"
)

//封包拆包类实例，暂时不需要成员
type DataPack struct{}

//封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func (this *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

//封包方法(压缩数据)
func (this *DataPack) Pack(msg iface.IMessagePackage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写Tag
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTag()); err != nil {
		return nil, err
	}

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法(解压数据)
func (this *DataPack) Unpack(binaryData []byte) (iface.IMessagePackage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &MessagePackage{}

	//读Tag
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Tag); err != nil {
		return nil, err
	}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if msg.Len > MAX_PKG_LEN {
		return nil, errors.New("Too large msg data recieved")
	}

	//TCP这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据,WebSocket全部读出
	if uint32(dataBuff.Len()) == msg.Len {
		//读data
		data := make([]byte, msg.Len)
		if _, err := dataBuff.Read(data); err != nil {
			return nil, err
		}
		msg.SetData(data)
	}

	return msg, nil
}
