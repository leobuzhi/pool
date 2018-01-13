package pool

type Pool interface {
	Get() interface{}
	Put()
	Len() int
	Close()
	Release()
}
