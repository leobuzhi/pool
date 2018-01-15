package main

import (
	"github.com/leobuzhi/pool"
	"time"
	"fmt"
	"net"
)

func init(){
	listerner,_:=net.Listen("tcp","127.0.0.1:8000")
	go listerner.Accept()
}

func main() {
	connFunc := func() (conn interface{}, err error) { return net.Dial("tcp", "127.0.0.1:8000") }
	closeFunc := func(conn interface{}) error { return conn.(net.Conn).Close() }

	config := pool.PoolConfig{
		InitConnNum: 1,
		MaxConnNum:  5,
		ConnFunc:    connFunc,
		CloseFunc:   closeFunc,
		IdelTime:    time.Second,
	}
	pool, err := pool.NewChannelPool(config)
	if err != nil {
		panic("new channelPool faild")
	}

	v, err := pool.Get()
	if err != nil {
		fmt.Println("pool get faild")
	}
	_, ok := v.(net.Conn)
	if !ok {
		fmt.Println("conn is not a net.Conn")
	}

	err = pool.Put(v)
	if err != nil {
		fmt.Println("pool put faild")
	}

	fmt.Println("pool len :", pool.Len())

	pool.Release()
}
