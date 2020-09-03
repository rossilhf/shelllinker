package update

//package main

import (
	"os/exec"
	//"runtime"
	"strings"
	"fmt"
	//"container/list"
	//"strconv"
)

//process device-part version update
func Update(url string) {
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
}

/*func main() {
	mac, ip := Get_curNet("linux")
	fmt.Println("result:")
	fmt.Println(mac)
	fmt.Println(ip)
}*/
