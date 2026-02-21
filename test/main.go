package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker        = "tcp://localhost:1883"
	topic         = "load/test"
	numClients    = 500
	publishPeriod = 2 * time.Second
)

func createClient(clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting %s: %v", clientID, token.Error())
	}
	return client
}

func main() {
	var wg sync.WaitGroup

	// Create Subscribers
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			clientID := fmt.Sprintf("sub-%d", i)
			client := createClient(clientID)

			client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
				// Uncomment if you want to see messages (very noisy)
				// fmt.Printf("[%s] %s\n", clientID, string(msg.Payload()))
			})
		}(i)
	}

	// Create Publishers
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			clientID := fmt.Sprintf("pub-%d", i)
			client := createClient(clientID)

			ticker := time.NewTicker(publishPeriod)
			defer ticker.Stop()

			for range ticker.C {
				payload := fmt.Sprintf("message from %s at %s",
					clientID, time.Now().Format(time.RFC3339))

				token := client.Publish(topic, 0, false, payload)
				token.Wait()
			}
		}(i)
	}

	wg.Wait()
}