package lnet

import (
	"github.com/golang/protobuf/proto"
)

type PbProtocol struct {

}

func (this *PbProtocol) Marshal(msg interface{})([]byte,error){
	pb, ok := msg.(proto.Message)
	if !ok {
		return nil,NewError("proto类型错误",0)
	}

	data,err := proto.Marshal(pb)
	if err != nil{
		return nil, NewError("proto编码失败",0)
	}

	return data, nil
}

func (this *PbProtocol) Unmarshal(data []byte, v interface{}) error{
	pb, ok := v.(proto.Message)
	if !ok {
		return NewError("proto类型错误",0)
	}

	return proto.Unmarshal(data, pb)
}

