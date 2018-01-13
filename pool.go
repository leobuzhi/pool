package pool

type Pool interface {
	Get() (interface{}, error)
	Put()
	Len() int
	Close()
	Release()
}
