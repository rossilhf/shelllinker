package main

import (
	"context"
	"fmt"
	//"os"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"time"

	"device/config"
	"device/encryption"
	"device/getinfo"
)

var (
	IOTREPORT = "iotreport"
	IOTRETURN = "iotreturn"
	IOTDEV    = "iot/dev"
)

func iotReportMsgHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Println("got:", msg.Topic(), string(msg.Payload()))
	retstr := "recvedsssssssssssssssddddddddddddddddddddd" + string(msg.Payload())
	client.Publish(IOTRETURN, 1, false, retstr)
}

func onConnectHandler(client MQTT.Client) {
	fmt.Println("client:", "connected")
	client.Subscribe("iot/+", 1, iotReportMsgHandler)
}

func publishTimer(ctx context.Context, client MQTT.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		strsend := strconv.Itoa(int(time.Now().Unix())) + "report__0"
		client.Publish(IOTREPORT, 1, false, strsend)
		time.Sleep(time.Second * 10)
	}
}

func main() {
	ip, user, psw := config.ReadMQinfo()
	fmt.Println(ip, user, psw)

	str := encryption.Encrypt(11, "tcp://117.50.109.189:1883")
	fmt.Println(str)
	str = encryption.Decrypt(11, str)
	fmt.Println(str)

	str = getinfo.Get_curVersion()
	fmt.Println(str)

	str = getinfo.Get_curOs()
	fmt.Println(str)

	//publish.py
	/*opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://117.50.109.189:1883")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	opts.SetOnConnectHandler(onConnectHandler)

	client := MQTT.NewClient(opts)
	client.Connect()
	c, cancel := context.WithCancel(context.Background())
	go publishTimer(c, client)

	sigChan := make(chan os.Signal)
	<-sigChan
	cancel()
	fmt.Println("shutting down server")*/
}
