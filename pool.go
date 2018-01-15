package pool

type Pool interface {
	Get() (interface{}, error)
	Put(interface{}) error
	Len() int
	Close(interface{}) error
	Release()
}
