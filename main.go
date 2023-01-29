package main

import (
	"fmt"

	"os/exec"
)

func main() {
	s, err := exec.LookPath("go")
	if err != nil {

	}
	fmt.Println(s)
}
