package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("pub-client")
	opts.SetCleanSession(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected (publisher)")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Println("Connection lost:", err)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for i := 0; i < 5; i++ {
		payload := fmt.Sprintf("hello %d", i)

		token := client.Publish("test/topic", 0, false, payload)
		token.Wait()

		log.Println("Published:", payload)
		time.Sleep(time.Second)
	}

	client.Disconnect(250)
}
