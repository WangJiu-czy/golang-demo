package ants

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

import "testing"

var runTimes = 1000

// 使用协程池
func TestCommonPool(t *testing.T) {
	//关闭协程池
	defer ants.Release()

	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}

	// 执行100次调用
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		err := ants.Submit(syncCalculateSum)
		if err != nil {
			panic(err)
		}
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")
}

// 使用带函数的协程池
func TestPoolWithFunc(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
		wg.Done()
	})

	defer p.Release()
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		//调用协程池中定义的方法
		_ = p.Invoke(int32(i))
	}
	wg.Wait()

	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}

func TestWithOptions(t *testing.T) {
	p, _ := ants.NewPool(3, ants.WithExpiryDuration(time.Duration(time.Second)))
	defer p.Release()
	for i := 0; i < 10; i++ {
		err := p.Submit(func() {
			fmt.Printf("执行方法....... %d \n ", i)
			time.Sleep(1 * time.Second) // 模拟耗时操作
			fmt.Println("Hello With Options")
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(5 * time.Second)
}
