package rpc

import (
	"errors"
	"net"
	"reflect"
)

type Client struct {
	conn  net.Conn
	coder Coder
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:  conn,
		coder: &DefaultCoder{},
	}

}

func (c *Client) CallRPC(rpcName string, fPtr interface{}) {
	container := reflect.ValueOf(fPtr).Elem()

	// 回调方法
	f := func(req []reflect.Value) []reflect.Value {
		cReqTransport := NewTransport(c.conn)
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		// Process input parameters
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}
		// ReqRPC
		reqRPC := RPCdata{Name: rpcName, Args: inArgs}
		b, err := c.coder.Encode(reqRPC)
		if err != nil {
			panic(err)
		}
		err = cReqTransport.Send(b)
		if err != nil {
			return errorHandler(err)
		}
		// receive response from server
		rsp, err := cReqTransport.Read()
		if err != nil { // local network error or decode error
			return errorHandler(err)
		}
		rspDecode, _ := c.coder.Decode(rsp)
		if rspDecode.Err != "" { // remote server error
			return errorHandler(errors.New(rspDecode.Err))
		}
		if len(rspDecode.Args) == 0 {
			rspDecode.Args = make([]interface{}, container.Type().NumOut())
		}
		// unpackage response arguments
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { // unpackage arguments (except error)
				if rspDecode.Args[i] == nil { // if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rspDecode.Args[i])
				}
			} else { // unpackage error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), f))
}
