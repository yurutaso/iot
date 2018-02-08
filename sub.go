package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
)

func main() {
	const MQTT_BROKER = "tcp://127.0.0.1:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("localhost")
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	var wg sync.WaitGroup
	wg.Add(1)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	token := client.Subscribe("go-mqtt/sample", 0, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("topic: %s\n", msg.Topic())
		fmt.Printf("message: %s\n", msg.Payload())
		wg.Done()
	})
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	if token := client.Publish("go-mqtt/sample", 0, false, `{"message": "hello"}`); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	wg.Wait()
}
