package lnet

import "unsafe"

type MyProtocol struct {

}

func (this *MyProtocol) Encode(msg *Message) *Pakge{
	pakge := &Pakge{}
	pakge.head.Tag = 11
	pakge.head.Len = uint16(unsafe.Sizeof(msg))
	pakge.data = make([]byte,pakge.head.Len)
	pakge.data = *(*[]byte)(unsafe.Pointer(msg))
	return pakge
}

func (this *MyProtocol) Decode(pakge *Pakge) *Message{
	msg := (*Message)(unsafe.Pointer(pakge))
	return msg
}

