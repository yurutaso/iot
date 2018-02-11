package main

import (
	"github.com/yurutaso/iot"
	"log"
)

func main() {
	broker := iot.NewBroker(`domain of broker`)
	broker.SetUserPassword(`username`, `password`)
	broker.SetClientID(`subscriber`) // any character is OK, but it must be different from other client's ID
	broker.SetCertFiles(
		`/path/to/ca.crt`,
		`/path/to/client.crt`,
		`/path/to/client.key`,
	)
	err := broker.Subscribe(`topic`, iot.PrintTopicMessage)
	if err != nil {
		log.Fatal(err)
	}
}
