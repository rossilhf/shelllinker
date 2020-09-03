package main

import (
	"context"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/exec"
	//"strconv"
	"encoding/json"
	"strings"
	"time"

	"shlkr_device/config"
	"shlkr_device/encryption"
	"shlkr_device/getinfo"
	"shlkr_device/update"
)

var cur_toolaccount string

//structure: device report to server info
type ReportDeviceInfoStru struct {
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
//type: exec:  exec cmd on device
//type: update: device-part version update
type ReceiveServerCmdStru struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
}

//structure: device report cmd result
type ReportCmdResultStru struct {
	Msgtype        string `json:"msgtype"`
	Curuser        string `json:"curuser"`
	Curtoolaccount string `json:"curtoolaccount"`
	Curmac         string `json:"curmac"`
	Curpath        string `json:"curpath"`
	Report         string `json:"report"`
}

//process device-part version update
/*func update(url string) {
	//download device bin file
	ok := true
	excute := exec.Command("/bin/sh", "-c", "wget "+url)
	_, err := excute.Output()
	if err != nil {
		fmt.Println(err)
		cmdresult := "Error: wget " + url + " err!"
		fmt.Println(cmdresult)
		ok = false
	}

	//remove old bin file
	if ok {
		excute := exec.Command("/bin/sh", "-c", "rm shlkr")
		_, err := excute.Output()
		if err != nil {
			ok = false
			fmt.Println(err)
			cmdresult := "Error: rm shlkr err!"
			fmt.Println(cmdresult)
		}
	}

	//change downloaded file name
	if ok {
		tmplist := strings.Split(url, "/")
		filename := tmplist[len(tmplist)-1]
		fmt.Println("filename: ", filename)
		excute := exec.Command("/bin/sh", "-c", "mv "+filename+" shlkr")
		_, err := excute.Output()
		if err != nil {
			ok = false
			fmt.Println(err)
			cmdresult := "Error: rm shlkr err!"
			fmt.Println(cmdresult)
		}
	}

	//chmod +x
	if ok {
		excute := exec.Command("/bin/sh", "-c", "chmod +x shlkr")
		_, err := excute.Output()
		if err != nil {
			ok = false
			fmt.Println(err)
			cmdresult := "Error: chmod +x shlkr err!"
			fmt.Println(cmdresult)
		}
	}

	//kill current ./shlkr process, restart new version ./shlkr process
	//note: auto-restart new process will by shlkr.sh
	if ok {
		excute := exec.Command("/bin/sh", "-c", "ps -ef | grep ./shlkr | grep -v grep | awk '{print $2}'|xargs kill -9")
		_, err := excute.Output()
		if err != nil {
			ok = false
			fmt.Println(err)
			cmdresult := "Error: kill old shlkr err!"
			fmt.Println(cmdresult)
		}
	}
}*/

//excute cmd from shelllinker server, e.g: ls, pwd, lsusb
//and report results
func execCmdHandler(client MQTT.Client, msg MQTT.Message) {
	//get cmd content
	cmd_encry := string(msg.Payload())
	fmt.Println("got from topic:", msg.Topic(), ", cmd:", cmd_encry)
	cmd := encryption.Decrypt(11, cmd_encry)
	fmt.Println("got from topic:", msg.Topic(), ", cmd(decrypt) ", cmd)

	cmdstru := ReceiveServerCmdStru{}
	_ = json.Unmarshal([]byte(cmd), &cmdstru)
	cmdtype := cmdstru.Type
	cmd = cmdstru.Cmd
	fmt.Println("got cmd type: ", cmdtype)
	fmt.Println("got cmd: ", cmd)

	if cmdtype == "update" {
		update.Update(cmd)
	}

	if cmdtype == "exec" {
		//excute cmd
		cmdresult := ""
		if strings.IndexAny(cmd, "cd ") == 0 {
			leng := len(cmd)
			//path := cmd[3 : leng-1]
			path := cmd[3:leng]
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

		cmdresultstru := ReportCmdResultStru{
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
}

//listen cmd from shelllinker server
func listenCmdHandler(client MQTT.Client) {
	fmt.Println("client:", "connected")

	curos := getinfo.Get_curOs()
	curmac, _ := getinfo.Get_curNet(curos)
	topic := "topic_ser2dev/exec_cmd/" + curmac
	client.Subscribe(topic, 1, execCmdHandler)
}

//report current device info every hour
func reportDevInfo(ctx context.Context, client MQTT.Client) {
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
		fmt.Println(cur_mac, cur_ip)

		reportinfo := ReportDeviceInfoStru{
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
		time.Sleep(time.Second * 60 * 5) //60)
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

	//connect to mqtt server, listen and execute cmd from shelllinker server
	opts := MQTT.NewClientOptions()
	opts.AddBroker(de_mqttip)
	opts.SetUsername(de_mqttuser)
	opts.SetPassword(de_mqttpsw)
	opts.SetOnConnectHandler(listenCmdHandler)

	//report device info every hour
	client := MQTT.NewClient(opts)
	client.Connect()
	c, cancel := context.WithCancel(context.Background())
	go reportDevInfo(c, client)

	//listen ctrl+c to exit
	sigChan := make(chan os.Signal)
	<-sigChan
	cancel()
	fmt.Println("shutting down server")
}
