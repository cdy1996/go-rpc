package rpc

import (
	"fmt"
	"log"
	"net"
	"reflect"
)

// RPCServer ...
type RPCServer struct {
	addr   string
	funcs  map[string]reflect.Value
	listen net.Listener
	close  bool
}

func NewServer(addr string) *RPCServer {
	server := &RPCServer{
		addr:  addr,
		funcs: make(map[string]reflect.Value),
		close: false,
	}
	listener, e := net.Listen("tcp", server.addr)
	if e != nil {
		panic(e)
	}
	server.listen = listener
	return server
}

func (self *RPCServer) Start() {
	for {
		conn, e := self.listen.Accept()
		if e != nil {
			panic(e)
		}
		go self.Response(conn)
	}
}

func (self *RPCServer) Response(conn net.Conn) {
	transport := Transport{conn: conn}
	for {
		read, e := transport.Read()
		if e != nil {
			return
		}
		cdata, e := Decode(read)
		if e != nil {
			panic(e)
		}
		result := self.Execute(cdata)

		encode, e := Encode(result)
		if e != nil {
			return
		}
		_ = transport.Send(encode)
	}
}

func (self *RPCServer) Close() {
	self.close = true
	e := self.listen.Close()
	if e != nil {
		panic(e)
	}
}

// Register the name of the function and its entries
func (s *RPCServer) Register(fnName string, fFunc interface{}) {
	if _, ok := s.funcs[fnName]; ok {
		return
	}
	s.funcs[fnName] = reflect.ValueOf(fFunc)
}

// Execute the given function if present
func (s *RPCServer) Execute(req RPCdata) RPCdata {
	// 获取方法名
	f, ok := s.funcs[req.Name]
	if !ok {
		//  方法不存在
		e := fmt.Sprintf("func %s not Registered", req.Name)
		log.Println(e)
		return RPCdata{Name: req.Name, Args: nil, Err: e}
	}
	log.Printf("func %s is called\n", req.Name)
	// 参数解析
	inArgs := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		inArgs[i] = reflect.ValueOf(req.Args[i])
	}
	// 反射调用
	out := f.Call(inArgs)

	// 返回结果 最后一个是err
	resArgs := make([]interface{}, len(out)-1)
	for i := 0; i < len(out)-1; i++ {
		// Interface returns the constant value stored in v as an interface{}.
		resArgs[i] = out[i].Interface()
	}
	// 单独将err取出
	var er string
	if e, ok := out[len(out)-1].Interface().(error); ok {
		// convert the error into error string value
		er = e.Error()
	}
	return RPCdata{Name: req.Name, Args: resArgs, Err: er}
}
