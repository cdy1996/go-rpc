package rpc

// RPC数据传输格式
type RPCdata struct {
	Name string        // name of the function
	Args []interface{} // request's or response's body expect error.
	Err  string        // Error any executing remote server
}
