package getinfo

//package main

import (
	"os/exec"
	"runtime"
	"strings"
	//"fmt"
	"container/list"
	"strconv"
)

func Get_curVersion() string {
	version := "V1.0.1"

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

type PackStru struct {
	PackNum int64
	PackIdx int
}

func Sort(oldList *list.List) (newlist *list.List) {
	newList := list.New()

	for v := oldList.Front(); v != nil; v = v.Next() {
		node := newList.Front()
		for nil != node {
			if node.Value.(PackStru).PackNum < v.Value.(PackStru).PackNum {
				newList.InsertBefore(v.Value.(PackStru), node)
				break
			}
			node = node.Next()
		}

		if node == nil {
			newList.PushBack(v.Value.(PackStru))
		}
	}

	return newList
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
		//asume max packets matches true net-card
		lineNum := len(curNetInfoSlice)
		maxPackIdx := 0
		list_pack := list.New()
		for i := 0; i < lineNum; i++ {
			if strings.Contains(curNetInfoSlice[i], "RX packets ") {
				rxlist := strings.Split(curNetInfoSlice[i], " ")   //example: [8_space]RX packets 285275080  bytes 27883859623 (27.8 GB)
				txlist := strings.Split(curNetInfoSlice[i+2], " ") //example: [8_space]TX packets 437094372  bytes 44359708792 (44.3 GB)
				rxpack, _ := strconv.ParseInt(rxlist[8+2], 10, 64)
				txpack, _ := strconv.ParseInt(txlist[8+2], 10, 64)
				packnum := rxpack + txpack
				list_pack.PushBack(PackStru{packnum, i})
			}
		}
		//for v := list_pack.Front(); v != nil; v = v.Next() {
		//	maxPackNum := v.Value.(PackStru).PackNum
		//	maxPackIdx := v.Value.(PackStru).PackIdx
		//	fmt.Println(maxPackNum)
		//	fmt.Println(maxPackIdx)
		//}
		//fmt.Println("")
		//fmt.Println("")
		new_list_pack := Sort(list_pack)
		//for v := new_list_pack.Front(); v != nil; v = v.Next() {
		//	maxPackNum := v.Value.(PackStru).PackNum
		//	maxPackIdx := v.Value.(PackStru).PackIdx
		//	fmt.Println(maxPackNum)
		//	fmt.Println(maxPackIdx)
		//}
		//fmt.Println("")
		//fmt.Println("")

		foundflag := false
		for v := new_list_pack.Front(); v != nil; v = v.Next() {
			if foundflag == true {
				break
			}
			maxPackIdx = v.Value.(PackStru).PackIdx

			//find mac address
			tmp := maxPackIdx - 15
			if tmp < 0 {
				tmp = 0
			}
			for i := maxPackIdx; i > tmp; i-- {
				if curNetInfoSlice[i] == "" {
					break
				}
				if strings.Contains(curNetInfoSlice[i], "ether ") {
					//fmt.Println(curNetInfoSlice[i])
					maclist := strings.Split(curNetInfoSlice[i], " ") //example: [8_space]ether 38:d5:47:00:42:52  txqueuelen 1000  (以太网)
					macAddr = maclist[8+1]
				}
			}

			//find ip address
			for i := maxPackIdx; i > tmp; i-- {
				//fmt.Println("kkkkkkkkkkk"+strconv.Itoa(i))
				if curNetInfoSlice[i] == "" {
					break
				}
				if strings.Contains(curNetInfoSlice[i], "inet ") {
					iplist := strings.Split(curNetInfoSlice[i], " ") //example: [8_space]inet 10.39.251.182  netmask 255.255.252.0  broadcast 10.39.251.255
					ipAddr = iplist[8+1]
				}
			}
			if ipAddr == "127.0.0.1" || ipAddr == "0.0.0.0" || ipAddr == "unknown" {
				//fmt.Println("found err ip:"+ipAddr)
				continue
			} else {
				//fmt.Println("found ip:"+ipAddr)
				foundflag = true
			}
			//fmt.Println("ssssssssssssssssssss"+ipAddr)
			//fmt.Println("dddddddddddddddddddd"+macAddr)
			//fmt.Println("uuuuuuuuuuuuuuuuuuuu"+strconv.Itoa(maxPackIdx))
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
	fmt.Println("result:")
	fmt.Println(mac)
	fmt.Println(ip)
}*/
