package main

import (
	"consensus/process"
	"fmt"
	"strconv"
)

var n int

func main() {
	fmt.Print("Hello. How many processes do you want?")
	fmt.Scanln(&n)
	var ports []string
	var entry = "localhost:"

	for i := range n {
		var number = i + 50051
		var port = entry + strconv.Itoa(number)
		ports = append(ports, port)
	}

	for i := range n {
		var port = ports[i]
		var id = i
		process.Run(port, id)

	}
}
