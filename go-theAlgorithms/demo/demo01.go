package demo

import (
	"fmt"
	"sync"
)

/*使⽤两个 goroutine 交替打印序列，⼀个 goroutine 打印数字， 另外⼀个 goroutine 打印字⺟， 最终效果
如下：
12AB34CD56EF78GH910IJ1112KL1314MN1516OP1718QR1920ST2122UV2324WX2526YZ2728
*/

/*
问题很简单，使⽤ channel 来控制打印的进度。使⽤两个 channel ，来分别控制数字和字⺟的打印序列， 数字打
印完成后通过 channel 通知字⺟打印, 字⺟打印完成后通知数字打印，然后周⽽复始的⼯作。
*/
func process() {
	letter, number := make(chan bool), make(chan bool)
	wait := sync.WaitGroup{}

	go func() {
		i := 1
		for {
			select {
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++

				letter <- true
			}
		}
	}()
	wait.Add(1)
	go func(wait *sync.WaitGroup) {
		i := 'A'
		for {
			select {
			case <-letter:
				if i >= 'Z' {
					wait.Done()

					return
				}
				fmt.Print(string(i))
				i++
				fmt.Print(string(i))
				i++
				number <- true
			}

		}

	}(&wait)
	number <- true
	wait.Wait()
}
