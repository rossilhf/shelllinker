package getinfo
//package main

import (
	"os/exec"
	"runtime"
	"strings"
	//"fmt"
	"strconv"
)

func Get_curVersion() string {
	version := "V20200831"

	return version
}

func Get_curOs() string {
	sysType := runtime.GOOS
	if (sysType == "linux") || (sysType == "windows") {
		return sysType
	} else {
		sysType = "unknown"
	}

	return sysType
}

func Get_curCpu(cur_os string) string {
	curCpu := "unknown"
	if cur_os == "linux" {
		cmd := exec.Command("lscpu")
		buf, _ := cmd.Output()
		curCpuInfo := string(buf)
		cpuinfoslice := strings.Split(curCpuInfo, "\n")
		curCpuInfo = cpuinfoslice[0]
		cpuinfoslice = strings.Split(curCpuInfo, " ")
		curCpu = cpuinfoslice[len(cpuinfoslice)-1]
	}

	if cur_os == "windows" {
		cmd := exec.Command("wmic cpu list brief")
		buf, _ := cmd.Output()
		curCpuInfo := string(buf)
		cpuinfoslice := strings.Split(curCpuInfo, "\n")
		curCpuInfo = cpuinfoslice[1]
		cpuinfoslice = strings.Split(curCpuInfo, " ")
		curCpu = cpuinfoslice[0]
	}

	return curCpu
}

//get current mac address and ip
//find mac/ip by most activate RX/TX packets
func Get_curNet(cur_os string) (string, string) {
	macAddr := "unknown"
	ipAddr := "unknown"

	if cur_os == "linux" {
		cmd := exec.Command("ifconfig")
		buf, _ := cmd.Output()
		curNetInfo := string(buf)
		curNetInfoSlice := strings.Split(curNetInfo, "\n")

		//find the line just like: "[8_space]RX packets 285275080  bytes 27883859623 (27.8 GB)"
		lineNum := len(curNetInfoSlice)
		var maxPackNum int64 = 0
		maxPackIdx := 0
		for i := 0; i < lineNum; i++ {
			if strings.Contains(curNetInfoSlice[i], "RX packets ") {
				rxlist := strings.Split(curNetInfoSlice[i], " ")   //example: [8_space]RX packets 285275080  bytes 27883859623 (27.8 GB)
				txlist := strings.Split(curNetInfoSlice[i+2], " ") //example: [8_space]TX packets 437094372  bytes 44359708792 (44.3 GB)
				rxpack, _ := strconv.ParseInt(rxlist[8+2], 10, 64)
				txpack, _ := strconv.ParseInt(txlist[8+2], 10, 64)
				packnum := rxpack + txpack
				if packnum > maxPackNum {
					maxPackNum = packnum
					maxPackIdx = i
				}
			}
		}

		//find mac address
		for i := maxPackIdx; i > maxPackIdx-5; i-- {
			if strings.Contains(curNetInfoSlice[i], "ether ") {
				maclist := strings.Split(curNetInfoSlice[i], " ") //example: [8_space]ether 38:d5:47:00:42:52  txqueuelen 1000  (以太网)
				macAddr = maclist[8+1]
			}
		}

		//find ip address
		for i := maxPackIdx; i > maxPackIdx-5; i-- {
			if strings.Contains(curNetInfoSlice[i], "inet ") {
				iplist := strings.Split(curNetInfoSlice[i], " ") //example: [8_space]inet 10.39.251.182  netmask 255.255.252.0  broadcast 10.39.251.255
				ipAddr = iplist[8+1]
			}
		}
	}

	if cur_os == "windows" {
	}

	return macAddr, ipAddr
}

func Get_curUser(cur_os string) string {
	cmd := exec.Command("whoami")
	buf, _ := cmd.Output()
	curUser := string(buf)
	curUser = strings.Replace(curUser, "\n", "", -1)

	return curUser
}

func Get_curPath(cur_os string) string {
	curPath := "unknown"

	if cur_os == "linux" {
		cmd := exec.Command("pwd")
		buf, _ := cmd.Output()
		curPath = string(buf)
	}

	if cur_os == "windows" {
	}

	return curPath
}

/*func main() {
	mac, ip := Get_curNet("linux")
	fmt.Println(mac)
	fmt.Println(ip)
}*/
