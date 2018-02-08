package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"time"
)

func NewTLSConfig() *tls.Config {
	certs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(`mosquitto.org.crt`)
	if err != nil {
		log.Fatal(err)
	}
	certs.AppendCertsFromPEM(pem)
	return &tls.Config{
		RootCAs: certs,
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	tlsconfig := NewTLSConfig()

	//const MQTT_BROKER = "ssl://127.0.0.1:8883"
	const MQTT_BROKER = "ssl://test.mosquitto.org:8883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("subscriber")
	opts.SetTLSConfig(tlsconfig)
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	token := client.Subscribe("/test/tls", 0, f)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
