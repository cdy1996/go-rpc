package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type User struct {
	Name string
	Age  int
}

var userDB = map[int]User{
	1: User{"Ankur", 85},
	9: User{"Anand", 25},
	8: User{"Ankur Anand", 27},
}

func QueryUser(id int) (User, error) {
	if u, ok := userDB[id]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
	//u , err := QueryUser(8)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("name: %s, age: %d \n", u.Name, u.Age)

	//gomodtest_base.NewIntCollection("!23")

	// new Type needs to be registered
	gob.Register(User{})
	addr := "localhost:3212"
	srv := NewServer(addr)
	// start server
	srv.Register("QueryUser", QueryUser)
	go srv.Run()
	// wait for server to start.
	time.Sleep(1 * time.Second)
	// start client
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	cli := NewClient(conn)
	var Query func(int) (User, error)
	cli.callRPC("QueryUser", &Query)
	u, err := Query(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
	u2, err := Query(8)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
}
