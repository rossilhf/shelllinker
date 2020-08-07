package main

import (
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	//connection configure
	//broker := "tcp://10.39.251.182:1883"
	broker := "tcp://117.50.109.189:1883"
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("test")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	fmt.Println("connected.")

	receiveCount := 0
	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("test", byte(1), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for receiveCount < 5 {
		fmt.Println(receiveCount)
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}

	client.Disconnect(250)
	fmt.Println("Sample Subscriber Disconnected")
}
