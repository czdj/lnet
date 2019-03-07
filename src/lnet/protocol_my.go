package lnet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"unsafe"
)

type MyProtocol struct {

}

func (this *MyProtocol) Encode(tag uint16, msg interface{}) []byte{
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		fmt.Println("encode error:%v", err)
	}

	head := &PakgeHead{}
	head.Tag = tag
	head.Len = uint16(len(buf.Bytes()))

	data := make([]byte,unsafe.Sizeof(PakgeHead{}))
	ptr := (*PakgeHead)(unsafe.Pointer(&data[0]))
	ptr.Len = head.Len
	ptr.Tag = head.Tag
	data = append(data,buf.Bytes()...)

	return data
}

func (this *MyProtocol) Decode(tag uint16, data []byte) interface{}{
	buf := bytes.Buffer{}
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	msg := MsgTypeInfo.NewMsg(tag)
	if err := dec.Decode(msg); err != nil {
		fmt.Println("decode error:", err)
	}

	return msg
}

