package lprotocol

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"unsafe"
	"lnet"
)

type GobProtocol struct {

}

func (this *GobProtocol) Encode(tag uint16, msg interface{}) []byte{
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		fmt.Println("encode error:%v", err)
	}

	head := &lnet.PakgeHead{}
	head.Tag = tag
	head.Len = uint16(len(buf.Bytes()))

	data := make([]byte,unsafe.Sizeof(lnet.PakgeHead{}))
	ptr := (*lnet.PakgeHead)(unsafe.Pointer(&data[0]))
	ptr.Len = head.Len
	ptr.Tag = head.Tag
	data = append(data,buf.Bytes()...)

	return data
}

func (this *GobProtocol) Decode(tag uint16, data []byte) interface{}{
	buf := bytes.Buffer{}
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	msg := lnet.MsgTypeInfo.NewMsg(tag)
	if err := dec.Decode(msg); err != nil {
		fmt.Println("decode error:", err)
	}

	return msg
}

