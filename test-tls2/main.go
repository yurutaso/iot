package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	MQTT_BROKER string = "ssl://FQDN:8883"
	username    string = ""
	password    string = ""
	cafile      string = "/path/to/ca.crt"
	crtfile     string = "/path/to/client.crt"
	keyfile     string = "/path/to/client.key"
)

func NewTLSConfig(clientID string) *MQTT.ClientOptions {
	certs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(cafile)
	if err != nil {
		log.Fatal(err)
	}
	certs.AppendCertsFromPEM(pem)

	cert, err := tls.LoadX509KeyPair(crtfile, keyfile)
	if err != nil {
		log.Fatal(err)
	}

	tlsconfig := &tls.Config{
		RootCAs:      certs,
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    nil,
		Certificates: []tls.Certificate{cert},
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID(clientID)
	if username != "" {
		opts.SetUsername(username)
	}
	if password != "" {
		opts.SetPassword(password)
	}
	opts.SetTLSConfig(tlsconfig)

	return opts
}

var f MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal(`Usage: tlsconnect [pub/sub]`)
	}
	switch args[1] {
	case "sub":
		opts := NewTLSConfig("subsciber")
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
	case "pub":
		opts := NewTLSConfig("publisher")
		client := MQTT.NewClient(opts)
		defer client.Disconnect(250)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
		token := client.Publish("/test/tls", 0, false, `{"message": "hello"}`)
		token.Wait()
	}
}
