package track

import (
	"fmt"
	"time"
)

// 给函数包装一个计时器
func GetExecuteTimeWrapFunc(f func()) func() {
	return func() {
		begin := time.Now()
		defer func() {
			end := time.Now()
			fmt.Println("execute time is", end.Sub(begin).Seconds())
		}()

		f()
	}
}
