package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No arguments given. Try again :)")
		return
	}
	names := os.Args[1:]
	name := strings.Join(names, " ")
	fmt.Printf("Hello %s, Welcome to the jungle\n", name)
}
