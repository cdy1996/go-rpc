package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

// Transport will use TLV protocol
type Transport struct {
	conn net.Conn // Conn is a generic stream-oriented network connection.
}

// NewTransport creates a Transport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

// Send TLV data over the network
func (t *Transport) Send(data []byte) error {
	// we will need 4 more byte then the len of data
	// as TLV header is 4bytes and in this header
	// we will encode how much byte of data
	// we are sending for this request.
	buf := make([]byte, 4+len(data))
	// 设置定长
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// Read TLV sent over the wire
func (t *Transport) Read() ([]byte, error) {
	//读取定长
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	// 读取消息
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
