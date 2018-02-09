package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

var f MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	const MQTT_BROKER = "tcp://127.0.0.1:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	/*
		ID of subscribe client and publish client must be different,
		if connecting to the same broker.
		(e.g. when testing pub/sub in the local mosquitto MQTT server.)
	*/
	opts.SetClientID("subscriber")
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	token := client.Subscribe("go-mqtt/sample", 0, f)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
