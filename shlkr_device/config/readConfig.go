package config 

import (
	"fmt"
	"os"
	//"io"
	"bufio"
	"strings"
	//"container/list"
)

/*read rabbit mq info: 
1. mqtt-server ip&port, 
2. mqtt-server user, 
3. mqtt-server psw
4. shelllinker tool account*/
func ReadMQinfo() (string, string, string, string){
	fileName := "./config.dat"
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	var ip string
	var user string
	var psw string
	var toolaccount string
	buf := bufio.NewReader(fi)
	for i := 0; i < 4; i++ {
		line, _, _ := buf.ReadLine()
		content := string(line)
		content = strings.Replace(content, " ", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		content = strings.Replace(content, "\t", "", -1)
		fmt.Println(i, content)
		
		if i == 0 {
			ip = content
		}
		if i == 1 {
			user = content
		}
		if i == 2 {
			psw = content
		}
		if i == 3 {
			toolaccount = content
		}
	}

	return ip, user, psw, toolaccount
}
