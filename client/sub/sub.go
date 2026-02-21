package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:8883")
	opts.SetClientID("sub-client")
	opts.SetCleanSession(true)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
		log.Printf("[RECV] topic=%s payload=%s\n", m.Topic(), string(m.Payload()))
	})

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected (subscriber)")

		token := c.Subscribe("test/#", 0, nil)
		token.Wait()
		if token.Error() != nil {
			log.Println("Subscribe failed:", token.Error())
			return
		}

		log.Println("Subscribed to test/topic")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Println("Connection lost:", err)
	}

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	select {}
}
