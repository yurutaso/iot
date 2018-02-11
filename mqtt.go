package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"time"
)

const (
	DOMAIN   string = "domain.net"
	USERNAME string = "usrname for mqtt"
	PASSWORD string = "password for mqtt"
	CAFILE   string = "/etc/mosquitto/ca_certificates/ca.crt"
	CRTFILE  string = "/etc/mosquitto/certs-client/client.crt"
	KEYFILE  string = "/etc/mosquitto/certs-client/client.key"
)

func NewTLSConfig(clientID string) (*MQTT.ClientOptions, error) {
	certs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(CAFILE)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pem)

	cert, err := tls.LoadX509KeyPair(CRTFILE, KEYFILE)
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
	opts.AddBroker("ssl://" + DOMAIN + ":8883")
	opts.SetClientID(clientID)
	if USERNAME != "" {
		opts.SetUsername(USERNAME)
	}
	if PASSWORD != "" {
		opts.SetPassword(PASSWORD)
	}
	opts.SetTLSConfig(tlsconfig)

	return opts, nil
}

var f MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG: %s\n", message.Payload())
}

func Subscribe(f MQTT.MessageHandler) error {
	opts, err := NewTLSConfig("subscriber")
	if err != nil {
		return err
	}
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	token := client.Subscribe("/myMQTT", 0, f)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
func Publish(message string) error {
	opts, err := NewTLSConfig("publisher")
	if err != nil {
		return err
	}
	client := MQTT.NewClient(opts)
	defer client.Disconnect(250)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	token := client.Publish("/myMQTT", 0, false, message)
	token.Wait()
	return nil
}
