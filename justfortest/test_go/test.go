package main

import "fmt"

/*func*/
func max(num1 int, num2 int) int {
	result := 0
	if (num1 > num2){
		result = num1
	}else{
		result = num2
	}

	return result
}

func main() {
	//var i int
	i := 5
	//var f float64
	f := 6.0
	//var bb bool
	bb := false
	//var s string
	s := "ddddddddddd"
	fmt.Println(i, f, bb, s)

	var a int = 21
	var b int = 10
	if(a == b){
		fmt.Println("==")
	}else{
		fmt.Println("!=")
	}

	for j := 5; j <= 10; j++ {
		fmt.Println(j)
	}

	ret := max(a, b)
	fmt.Println("max num is: ", ret)
}
