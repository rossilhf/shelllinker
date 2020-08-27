package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	//"strconv"
	"time"

	"device/config"
	"device/encryption"
	"device/getinfo"
)

//var (
//	IOTREPORT = "iotreport"
//	IOTRETURN = "iotreturn"
//	IOTDEV    = "iot/dev"
//)

//excute cmd from shelllinker server, e.g: ls, pwd, lsusb
func iotReportMsgHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Println("got from topic:", msg.Topic(), " ", string(msg.Payload()))
	cmd := exec.Command(string(msg.Payload()))
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	cmdresult := string(buf)
	//fmt.Println("exec result: ", cmdresult)

	curos := getinfo.Get_curOs()
	curmac, _:= getinfo.Get_curNet(curos)
	topic := "topic_dev2ser/exec_result/" + curmac
	//client.Publish(IOTRETURN, 1, false, cmdresult)
	client.Publish(topic, 1, false, cmdresult)
}

//listen cmd from shelllinker server
func onConnectHandler(client MQTT.Client) {
	fmt.Println("client:", "connected")

	curos := getinfo.Get_curOs()
	curmac, _:= getinfo.Get_curNet(curos)
	topic := "topic_ser2dev/exec_cmd/" + curmac
	client.Subscribe(topic, 1, iotReportMsgHandler)
}

//heart beat
//report current device info every hour
func publishTimer(ctx context.Context, client MQTT.Client) { //, curmac string, curip string, curcpu string, curos string, curuser string, curversion string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		curversion := getinfo.Get_curVersion()
		fmt.Println(curversion)

		curos := getinfo.Get_curOs()
		fmt.Println(curos)

		curcpu := getinfo.Get_curCpu(curos)
		fmt.Println(curcpu)

		curmac, curip := getinfo.Get_curNet(curos)
		fmt.Println(curmac, curip)

		curuser := getinfo.Get_curUser(curos)
		fmt.Println(curuser)

		curpath := getinfo.Get_curPath(curos)
		fmt.Println(curpath)

		//strsend := strconv.Itoa(int(time.Now().Unix())) + "report__0"
		strsend := "{curmac:"+curmac + ", curip:"+curip + ", curcpu:"+curcpu + ", curos:"+curos + ", curuser:"+curuser + ", curversion:"+curversion+"}"
		strsend_encry := encryption.Encrypt(11, strsend)
		topic := "topic_dev2ser/dev_info/" + curmac
		client.Publish(topic, 1, false, strsend_encry)
		time.Sleep(time.Second * 10)
	}
}

func main() {
	ip, user, psw := config.ReadMQinfo()
	fmt.Println(ip, user, psw)

	//str := encryption.Encrypt(11, "tcp://117.50.109.189:1883")
	//fmt.Println(str)
	de_ip := encryption.Decrypt(11, ip)
	de_user := encryption.Decrypt(11, user)
	de_psw := encryption.Decrypt(11, psw)
	fmt.Println(de_ip, de_user, de_psw)

	//connect to mqtt server, execute cmd from shelllinker server
	opts := MQTT.NewClientOptions()
	opts.AddBroker(de_ip)
	opts.SetUsername(de_user)
	opts.SetPassword(de_psw)
	opts.SetOnConnectHandler(onConnectHandler)

	//report device info every hour
	client := MQTT.NewClient(opts)
	client.Connect()
	c, cancel := context.WithCancel(context.Background())
	go publishTimer(c, client)//, curmac, curip, curcpu, curos, curuser, curversion)
	
	//listen ctrl+c to exit
	sigChan := make(chan os.Signal)
	<-sigChan
	cancel()
	fmt.Println("shutting down server")
}
