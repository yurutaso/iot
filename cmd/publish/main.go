package main

import (
	"fmt"
	"github.com/yurutaso/iot"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: publish text")
		return
	}
	iot.Publish(args[1])
}
