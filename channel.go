package pool

import (
	"net"
	"time"
)

type PoolConfig struct {
	InitConnNum int
	MaxConnNum  int
	ConnFunc    func()
	CloseFunc   func(k interface{})
	IdelTime    time.Duration
}

type channelPool struct {
	pool       chan *idleConn
	maxConnNum int
	connFunc   func()
	closeFunc  func(k interface{})
	idelTime   time.Duration
}

type idleConn struct {
	conn     net.Conn
	idelTime time.Duration
}

func NewChannelPool(config PoolConfig) Pool {
	p := &channelPool{
		pool:       make(chan *idleConn, config.InitConnNum),
		maxConnNum: config.MaxConnNum,
		connFunc:   config.ConnFunc,
		closeFunc:  config.CloseFunc,
		idelTime:   config.IdelTime,
	}
	return p
}

func (pool channelPool) getConns() interface{} {
	return nil
}

func (pool channelPool) Get() interface{} {
	return nil
}

func (pool channelPool) Put() {

}

func (pool channelPool) Len() int {
	return 0
}

func (pool channelPool) Close() {

}

func (pool channelPool) Release() {

}
