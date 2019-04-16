package lnet

import (
	"github.com/golang/protobuf/proto"
	"unsafe"
)

type PbProtocol struct {

}

func (this *PbProtocol) Encode(tag uint16, msg interface{}) []byte{
	pb, ok := msg.(proto.Message)
	if !ok {
		return nil
	}

	encodeData,err := proto.Marshal(pb)
	if err != nil{
		return nil
	}

	head := &PakgeHead{}
	head.Tag = tag
	head.Len = uint16(len(encodeData))

	data := make([]byte,unsafe.Sizeof(PakgeHead{}))
	ptr := (*PakgeHead)(unsafe.Pointer(&data[0]))
	ptr.Len = head.Len
	ptr.Tag = head.Tag
	data = append(data,encodeData...)

	return data
}

func (this *PbProtocol) Decode(tag uint16, data []byte) interface{} {
	msg := MsgTypeInfo.NewMsg(tag)
	pb, ok := msg.(proto.Message)
	if !ok {
		return nil
	}
	proto.Unmarshal(data, pb)

	return pb
}

