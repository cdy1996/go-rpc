package rpc

import (
	"bytes"
	"encoding/gob"
)

type Coder interface {
	Encode(data RPCdata) ([]byte, error)
	Decode(b []byte) (RPCdata, error)
}

type DefaultCoder struct {
}

// 编码
func (self *DefaultCoder) Encode(data RPCdata) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//解码
func (self *DefaultCoder) Decode(b []byte) (RPCdata, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data RPCdata
	if err := decoder.Decode(&data); err != nil {
		return RPCdata{}, err
	}
	return data, nil
}
