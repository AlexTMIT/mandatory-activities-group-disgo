package main

import (
	"consensus/process"
	"fmt"
	"strconv"
)

var n int

func main() {
	fmt.Print("hello. how many processes do you want?")
	fmt.Scanln(&n)
	var ports []string
	var entry = "localhost:"

	for i := 0; i < n; i++ {
		var number = i + 50051
		var port = entry + strconv.Itoa(number)
		ports = append(ports, port)
	}

	for i := 0; i < n; i++ {
		var port = ports[i]
		go process.Run(port, ports)
	}

	// keep the main function running
	select {}
}
