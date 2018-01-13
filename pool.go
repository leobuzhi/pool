package pool

type Pool interface {
	Get() interface{}
	getConns() interface{}
	Put()
	Len() int
	Close()
	Release()
}
