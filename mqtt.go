package iot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"time"
)

type Broker struct {
	domain        string
	username      string
	password      string
	cafile        string
	clientCrtfile string
	clientKeyfile string
	clientID      string
}

func NewBroker(domain string) *Broker {
	return &Broker{domain: domain}
}

func (broker *Broker) SetUserPassword(username, password string) {
	broker.username = username
	broker.password = password
}

func (broker *Broker) SetClientID(clientID string) {
	broker.clientID = clientID
}

func (broker *Broker) SetCertFiles(cafile, crtfile, keyfile string) {
	broker.cafile = cafile
	broker.clientCrtfile = crtfile
	broker.clientKeyfile = keyfile
}

func (broker *Broker) NewTLSConfig() (*MQTT.ClientOptions, error) {
	certs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(broker.cafile)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pem)

	cert, err := tls.LoadX509KeyPair(broker.clientCrtfile, broker.clientKeyfile)
	if err != nil {
		return nil, err
	}

	tlsconfig := &tls.Config{
		RootCAs:      certs,
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    nil,
		Certificates: []tls.Certificate{cert},
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://" + broker.domain + ":8883")
	opts.SetClientID(broker.clientID)
	if broker.username != "" {
		opts.SetUsername(broker.username)
	}
	if broker.password != "" {
		opts.SetPassword(broker.password)
	}
	opts.SetTLSConfig(tlsconfig)

	return opts, nil
}

func PrintTopicMessage(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Topic: %s\n", message.Topic())
	fmt.Printf("Message: %s\n", message.Payload())
}

func (broker *Broker) Subscribe(topic string, handler MQTT.MessageHandler) error {
	opts, err := broker.NewTLSConfig()
	if err != nil {
		return err
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)
	token := client.Subscribe(topic, 0, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func (broker *Broker) Publish(topic, message string) error {
	opts, err := broker.NewTLSConfig()
	if err != nil {
		return err
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	return nil
}
