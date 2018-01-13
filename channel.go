package pool

import (
	"time"
	"sync"
	"fmt"
)

type PoolConfig struct {
	InitConnNum int
	MaxConnNum  int
	ConnFunc    func() (interface{}, error)
	CloseFunc   func(interface{}) error
	IdelTime    time.Duration
}

type channelPool struct {
	mu        sync.Mutex
	conns     chan *idleConn
	connFunc  func() (interface{}, error)
	closeFunc func(interface{}) error
	idelTime  time.Duration
}

type idleConn struct {
	conn     interface{}
	idelTime time.Time
}

func NewChannelPool(config PoolConfig) (Pool, error) {
	p := &channelPool{
		conns:     make(chan *idleConn, config.InitConnNum),
		connFunc:  config.ConnFunc,
		closeFunc: config.CloseFunc,
		idelTime:  config.IdelTime,
	}
	for i := 0; i != config.InitConnNum; i++ {
		conn, err := p.connFunc()
		if err != nil {
			p.closeFunc(conn)
			return nil, fmt.Errorf("create connection faild")
		}
		p.conns <- &idleConn{conn: conn, idelTime: time.Now()}
	}
	return p, nil
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
