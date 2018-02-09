package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func main() {
	//const MQTT_BROKER = "tcp://127.0.0.1:1883"
	const MQTT_BROKER = "tcp://192.168.11.40:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("publisher")
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	token := client.Publish("go-mqtt/sample", 0, false, `{"message": "hello"}`)
	token.Wait()
}
