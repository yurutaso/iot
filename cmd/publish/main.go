package main

import (
	"fmt"
	"github.com/yurutaso/iot"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: publish text")
		return
	}
	broker := iot.NewBroker(`domain of broker`)
	broker.SetUserPassword(`username`, `password`)
	broker.SetClientID(`publisher`) // any character is OK, but it must be different from other client's ID
	broker.SetCertFiles(
		`/path/to/ca.crt`,
		`/path/to/client.crt`,
		`/path/to/client.key`,
	)
	err := broker.Publish(`topic`, args[1])
	if err != nil {
		log.Fatal(err)
	}
}
