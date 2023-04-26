package main

import (
	"github.com/qianlnk/pgbar"
	"time"
)

func ProgressBar() {
	bar := pgbar.NewBar(0, "下载进度", 100)
	for i := 0; i < 100; i++ {
		bar.Add()
		time.Sleep(time.Millisecond * 30)
	}
}
