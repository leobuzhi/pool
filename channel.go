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
		conns:     make(chan *idleConn, config.MaxConnNum),
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

func (pool channelPool) getConns() chan *idleConn {
	pool.mu.Lock()
	conns := pool.conns
	pool.mu.Unlock()
	return conns
}

func (pool channelPool) Get() (interface{}, error) {
	conns := pool.getConns()
	if conns == nil {
		return nil, fmt.Errorf("get connections faild")
	}

	for {
		select {
		case wrapConn := <-conns:
			if timeOut := pool.idelTime; timeOut > 0 {
				if wrapConn.idelTime.Add(timeOut).Before(time.Now()) {
					pool.closeFunc(wrapConn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			conn, err := pool.connFunc()
			if err != nil {
				return nil, err
			}
			return conn, err
		}
	}
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
