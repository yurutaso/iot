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
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("samplecerts/CAfile.pem")
	if err != nil {
		log.Fatal(err)
	}
	certpool.AppendCertsFromPEM(pemCerts)
	cert, err := tls.LoadX509KeyPair("samplecerts/client-crt.pem", "samplecerts/client-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func main() {
	tlsconfig := NewTLSConfig()

	const MQTT_BROKER = "ssl://127.0.0.1:8883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTT_BROKER)
	opts.SetClientID("subscriber")
	opts.SetTLSConfig(tlsconfig)
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	token := client.Publish("go-mqtt/sample", 0, false, `{"message": "hello"}`)
	token.Wait()
}
