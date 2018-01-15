package pool

import (
	"time"
	"sync"
	"errors"
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
	if config.InitConnNum < 0 || config.InitConnNum > config.MaxConnNum || config.MaxConnNum <= 0 {
		return nil, errors.New("config error")
	}
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
			return nil, errors.New("create connection faild")
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
		return nil, errors.New("get connections faild")
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

func (pool channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("conn is nil")
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.conns == nil {
		return pool.closeFunc(conn)
	}

	select {
	case pool.conns <- &idleConn{conn: conn, idelTime: time.Now()}:
		return nil
	default:
		return pool.closeFunc(conn)
	}
}

func (pool channelPool) Len() int {
	return len(pool.conns)
}

func (pool channelPool) Close(conn interface{}) error {
	return pool.closeFunc(conn)
}

func (pool channelPool) Release() {
	pool.mu.Lock()
	conns := pool.conns
	closeFunc := pool.closeFunc
	pool.conns = nil
	pool.connFunc = nil
	pool.closeFunc = nil
	pool.mu.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for wrapConn := range conns {
		closeFunc(wrapConn.conn)
	}
}
