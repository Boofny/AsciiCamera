package main

import (
	"fmt"
	"time"
)

func main() {
	
	var s string
	var r string 
	for range 10{
		s += "time\n"
	}

	for range 10{
		r += "slime\n"
	}
	
	for range 10{
		time.Sleep(time.Second)
		fmt.Print(s)
		fmt.Printf("\033[%dA", 10)
		fmt.Print(r)
	}

	fmt.Println(s)
  fmt.Printf("\033[%dA", 10)
}
