package lprotocol

import (
	"bytes"
	"encoding/gob"
	"go.uber.org/zap"
	"lnet"
)

type GobProtocol struct {
}

func (this *GobProtocol) Marshal(msg interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		lnet.Logger.Error("gob encode error", zap.Any("err", err))
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *GobProtocol) Unmarshal(data []byte, v interface{}) error {
	buf := bytes.Buffer{}
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(v); err != nil {
		lnet.Logger.Error("gob decode error", zap.Any("err", err))
		return err
	}

	return nil
}
