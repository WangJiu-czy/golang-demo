package main

import (
	"flag"
	"fmt"
)

func main() {
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()
	if *showVersion {
		fmt.Println("66", *showVersion)
	}
	fmt.Println("11", showVersion)
}
