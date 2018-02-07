package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("topic: %s\ns", msg.Topic())
	fmt.Printf("message: %s\n", msg.Payload())
}

func main() {
	const MQTT_BROKER = "tcp://127.0.0.1:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("localhost")
	opts.SetDefaultPublishHandler(f)
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	time.Sleep(3 * time.Second)
	token := client.Subscribe("go-mqtt/sample", 0, f)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	for {
		time.Sleep(1 * time.Second)
	}
	defer client.Disconnect(250)
}
