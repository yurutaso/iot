package main

import (
	"github.com/yurutaso/iot"
)

func main() {
	broker := iot.NewBroker(`domain of mqtt broker`)
	broker.SetUserPassword(`username`, `password`)
	broker.SetClientID(`clientID`) // any characters is OK, but it must be different from other client's ID
	broker.SetCertFiles(
		`/path/to/ca.crt`,
		`/path/to/client.crt`,
		`/path/to/client.key`,
	)

	webhookhandler := iot.NewWebhookHandler(broker)
	webhookhandler.SetTopic(`topic`)

	domain := `domain of https server`
	iot.HttpsServer(domain, webhookhandler.PublishPost)
}
