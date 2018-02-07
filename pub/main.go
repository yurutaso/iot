package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func main() {
	const MQTT_BROKER = "tcp://127.0.0.1:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("localhost")
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	token := client.Publish("go-mqtt/sample", 0, false, `{"message": "hello"}`)
	token.Wait()
}
