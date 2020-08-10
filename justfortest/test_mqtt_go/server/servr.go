package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	username  = "admin"
	password  = "public"
	brokerUrl = "tcp://117.50.109.189:1883"
)

var (
	IOTREPORT = "iotreport"
	IOTRETURN = "iotreturn"
	IOTDEV    = "iot/dev"
)

func iotReturnHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Println("got: ", msg.Topic(), string(msg.Payload()))
}

func iotReportHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Println("got:", msg.Topic(), string(msg.Payload()))
}

func onConnectHandler(client MQTT.Client) {
	fmt.Println("client:", "connected")
	client.Subscribe(IOTREPORT, 1, iotReportHandler)
	client.Subscribe(IOTRETURN, 1, iotReturnHandler)
}

func main() {
	//publish.py
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerUrl)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetOnConnectHandler(onConnectHandler)

	client := MQTT.NewClient(opts)
	client.Connect()
	//c, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-sigChan:
			break
		default:
		}
		fmt.Println("cmd: ")
		var Input string
		fmt.Scanln(&Input)
		strsend := strconv.Itoa(int(time.Now().Unix())) + "---->" + Input
		client.Publish(IOTDEV, 1, false, strsend)
		fmt.Println("send cmd:", strsend)
	}

	fmt.Println("shutting down server")
}
