package main

import (
	"fmt"
	"os/exec"
)

func main() {
	/*fmt.Println("hello, world!")

	var a string = "rossi"
	fmt.Println(a)

	var b, c int = 1, 2
	fmt.Println(b, c)

	var i int
	var f float64
	var bb bool
	var s string
	fmt.Println(i, f, bb, s)*/

	excute := exec.Command("uname","-a")
	buf, err := excute.Output()
	if err != nil {
		fmt.Println(err)
	}
	cmdresult := string(buf)
	fmt.Println("exec result: ", cmdresult)
}
