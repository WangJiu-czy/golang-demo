package track

import (
	"fmt"
	"testing"
	"time"
)

func TestGetExecuteTimeWrapFunc(t *testing.T) {
	wrapFunc := GetExecuteTimeWrapFunc(do)
	wrapFunc()
}

func do() {
	fmt.Println("=============================")
	fmt.Println("进行逻辑处理............")
	time.Sleep(5 * time.Second)
	fmt.Println("=============================")
}
