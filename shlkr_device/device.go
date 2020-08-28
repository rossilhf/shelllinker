package main

import (
	"context"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/exec"
	//"strconv"
	"strings"
	"encoding/json"
	"time"

	"shlkr_device/config"
	"shlkr_device/encryption"
	"shlkr_device/getinfo"
)

var cur_toolaccount string

//structure: device report to server info
type ReportDeviceInfo struct {
	Msgtype        string `json:"msgtype"`
	Curtoolaccount string `json:"curtoolaccount"`
	Curmac         string `json:"curmac"`
	Curcpu         string `json:"curcpu"`
	Curos          string `json:"curos"`
	Curip          string `json:"curip"`
	Heartbeattime  string `json:"heartbeattime"`
	Curversion     string `json:"curversion"`
	Curuser        string `json:"curuser"`
}

//structure: device receive cmd from server
type ReceiveServerCmd struct {
	Cmd string `json:"cmd"`
}

//structure: device report cmd result
type ReportCmdResult struct {
	Msgtype        string `json:"msgtype"`
	Curuser        string `json:"curuser"`
	Curtoolaccount string `json:"curtoolaccount"`
	Curmac         string `json:"curmac"`
	Curpath        string `json:"curpath"`
	Report         string `json:"report"`
}

//excute cmd from shelllinker server, e.g: ls, pwd, lsusb
func iotReportMsgHandler(client MQTT.Client, msg MQTT.Message) {
	//get cmd content
	cmd_encry := string(msg.Payload())
	fmt.Println("got from topic:", msg.Topic(), ", cmd:", cmd_encry)
	cmd := encryption.Decrypt(11, cmd_encry)
	fmt.Println("got from topic:", msg.Topic(), ", cmd(decrypt) ", cmd)

	cmdstru := ReceiveServerCmd{}
	_ = json.Unmarshal([]byte(cmd), &cmdstru)
	cmd = cmdstru.Cmd
	fmt.Println("got cmd: ", cmd)

	//excute cmd
	cmdresult := ""
	if strings.IndexAny(cmd, "cd ") == 0 {
		leng := len(cmd)
		//path := cmd[3 : leng-1]
		path := cmd[3 : leng]
		err := os.Chdir(path)
		if err != nil {
			fmt.Println(err)
			cmdresult = "Error: excute " + cmd + " err!" 
		}
	} else {
		excute := exec.Command("/bin/sh", "-c", cmd)
		buf, err := excute.Output()
		if err != nil {
			fmt.Println(err)
			cmdresult = "Error: excute " + cmd + " err!" 
		} else {
			cmdresult = string(buf)
			fmt.Println("exec result: ", cmdresult)
		}
	}

	// generate result structure/json/encrypt-json
	cur_os := getinfo.Get_curOs()
	cur_path := getinfo.Get_curPath(cur_os)
	cur_mac, _ := getinfo.Get_curNet(cur_os)
	cur_user := getinfo.Get_curUser(cur_os)

	cmdresultstru := ReportCmdResult{
		Msgtype:        "exec_return",
		Curuser:        cur_user,
		Curtoolaccount: cur_toolaccount,
		Curmac:         cur_mac,
		Curpath:        cur_path,
		Report:         cmdresult,
	}
	jsonBytes, err := json.Marshal(cmdresultstru)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(jsonBytes))
	}
	cmdresult_encry := encryption.Encrypt(11, string(jsonBytes))

	//send out cmd excute result
	topic := "topic_dev2ser/exec_result/" + cur_mac
	client.Publish(topic, 1, false, cmdresult_encry)
}

//listen cmd from shelllinker server
func onConnectHandler(client MQTT.Client) {
	fmt.Println("client:", "connected")

	curos := getinfo.Get_curOs()
	curmac, _ := getinfo.Get_curNet(curos)
	topic := "topic_ser2dev/exec_cmd/" + curmac
	client.Subscribe(topic, 1, iotReportMsgHandler)
}

//report current device info every hour
func publishTimer(ctx context.Context, client MQTT.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		cur_version := getinfo.Get_curVersion()
		cur_os := getinfo.Get_curOs()
		cur_cpu := getinfo.Get_curCpu(cur_os)
		cur_mac, cur_ip := getinfo.Get_curNet(cur_os)
		cur_user := getinfo.Get_curUser(cur_os)

		reportinfo := ReportDeviceInfo{
			Msgtype:        "info_report",
			Curtoolaccount: cur_toolaccount,
			Curmac:         cur_mac,
			Curcpu:         cur_cpu,
			Curos:          cur_os,
			Curip:          cur_ip,
			Heartbeattime:  time.Now().Format("2006-01-02 15:04:05"),
			Curversion:     cur_version,
			Curuser:        cur_user,
		}
		//strsend := strconv.Itoa(int(time.Now().Unix())) + "report__0"
		jsonBytes, err := json.Marshal(reportinfo)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(jsonBytes))
		}
		strsend_encry := encryption.Encrypt(11, string(jsonBytes))
		topic := "topic_dev2ser/dev_info/" + cur_mac
		client.Publish(topic, 1, false, strsend_encry)
		time.Sleep(time.Second * 60 * 60)
	}
}

func main() {
	ip, user, psw, toolaccount := config.ReadMQinfo()
	fmt.Println(ip, user, psw, toolaccount)

	//str := encryption.Encrypt(11, "tcp://117.50.109.189:1883")
	//fmt.Println(str)
	de_mqttip := encryption.Decrypt(11, ip)
	de_mqttuser := encryption.Decrypt(11, user)
	de_mqttpsw := encryption.Decrypt(11, psw)
	de_toolaccount := encryption.Decrypt(11, toolaccount)
	fmt.Println(de_mqttip, de_mqttuser, de_mqttpsw, de_toolaccount)
	cur_toolaccount = de_toolaccount

	//connect to mqtt server, execute cmd from shelllinker server
	opts := MQTT.NewClientOptions()
	opts.AddBroker(de_mqttip)
	opts.SetUsername(de_mqttuser)
	opts.SetPassword(de_mqttpsw)
	opts.SetOnConnectHandler(onConnectHandler)

	//report device info every hour
	client := MQTT.NewClient(opts)
	client.Connect()
	c, cancel := context.WithCancel(context.Background())
	go publishTimer(c, client)

	//listen ctrl+c to exit
	sigChan := make(chan os.Signal)
	<-sigChan
	cancel()
	fmt.Println("shutting down server")
}
