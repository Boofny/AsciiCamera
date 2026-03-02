package main

import (
	"fmt"

	"golang.org/x/term"
)

func main(){

	w, h, err := term.GetSize(0)
	if err != nil {
		return 
	}
	fmt.Print(w, h)
}
