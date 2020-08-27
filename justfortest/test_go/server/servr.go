package main

import (
	"fmt"
	//"strconv"
	//"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	username  = "admin"
	password  = "public"
	brokerUrl = "tcp://117.50.109.189:1883"
)

var (
	IOTREPORT = "topic_dev2ser/dev_info/38:d5:47:00:42:52" //"iotreport"
	IOTRETURN = "topic_dev2ser/exec_result/38:d5:47:00:42:52" //"iotreturn"
	IOTDEV    = "topic_ser2dev/exec_cmd/38:d5:47:00:42:52" //"iot/dev"
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

	//sigChan := make(chan os.Signal)
	//signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT,
	//	syscall.SIGQUIT)
	//go func() {
	//	for s := range sigChan {
	//		switch s {
	//		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	//			fmt.Println("Program Exit...", s)
	//			os.Exit(1)
	//		default:
	//
	//		}
	//	}
	//}()

	readloop:
	for {
		//select {
		//case <-sigChan:
		//	return
		//default:
		//
		//}
		fmt.Println("cmd: ")
		var Input string
		_,err:=fmt.Scanln(&Input)
		if Input == "exit" || err != nil{
			break readloop
		}
		//strsend := strconv.Itoa(int(time.Now().Unix())) + "---->" + Input
		strsend := Input
		client.Publish(IOTDEV, 1, false, strsend)
		fmt.Println("send cmd:", strsend)
	}

	fmt.Println("shutting down server")
}
