package config 

import (
	"fmt"
	"os"
	//"io"
	"bufio"
	"strings"
	//"container/list"
)

/*read rabbit mq info: ip, user, psw*/
func ReadMQinfo() (string, string, string){
	fileName := "./configfiles/mqserver.dat"
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	var ip string
	var user string
	var psw string
	buf := bufio.NewReader(fi)
	for i := 0; i < 3; i++ {
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
	}

	return ip, user, psw
}
