package singleton

import "sync"

type singleton struct {
}

var (
	ins  *singleton
	once sync.Once
)

func GetIns() *singleton {
	once.Do(func() {
		ins = &singleton{}
	})
	return ins
}
