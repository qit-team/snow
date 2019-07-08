package close

import "sync"

var (
	closeSet []Closeable
	lock     sync.RWMutex
)

type Closeable interface {
	Close() (error)
}

//注册应用停止时需要释放链接的服务
func Register(closeable Closeable) {
	lock.Lock()
	defer lock.Unlock()
	closeSet = append(closeSet, closeable)
}

//批量注册应用停止时需要释放链接的服务
func MultiRegister(closeableSet ...Closeable) {
	lock.Lock()
	defer lock.Unlock()
	closeSet = append(closeSet, closeableSet...)
}

//释放链接
func Free() {
	for _, v := range closeSet {
		if v != nil {
			v.Close()
		}
	}
}
