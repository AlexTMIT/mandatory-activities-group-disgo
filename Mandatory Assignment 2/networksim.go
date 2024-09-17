package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var amount = 0
var finished bool

type client struct {
	seq   int
	ready bool
	ack   int
	ch    chan int
}

type server struct {
	seq   int
	ready bool
	ack   int
	ch    chan int
}

func main() {
	wg := new(sync.WaitGroup)
	initThread(wg)
	wg.Wait()
}

func initThread(wg *sync.WaitGroup) {
	wg.Add(2)

	client := &client{
		seq: rand.Intn(100),
		ch:  make(chan int, 2),
		ack: 0,
	}

	server := &server{
		seq: rand.Intn(100),
		ch:  make(chan int, 2),
		ack: 0,
	}

	go func() {
		defer wg.Done()
		clientGo(client, server)
	}()
	go func() {
		defer wg.Done()
		serverGo(client, server)
	}()
}

func clientGo(client *client, server *server) {
	fmt.Println("running client")
	for !finished {
		fmt.Println("running client in forloop")
		if !client.ready {
			server.ch <- client.seq
			server.ready = true
		} else {
			fmt.Println("running client in else statement")
			client.seq = <-client.ch //server's ack no.
			<-client.ch

			var msg = <-client.ch
			client.ack = msg + 1

			server.ch <- client.seq
			server.ch <- client.ack
			finished = true
			break
		}
	}
	fmt.Println("Connection established!")
	// client vil sende seq nummer ind til server
	// client skal herefter vente til, at client får 2 numre ind i sin egen kanal
	// når client får besken, vil den gerne sende serveren's seq + 1, plus sin egen seq nummer + 1
}

func serverGo(client *client, server *server) {
	fmt.Println("running server")
	for server.ready {
		fmt.Println("running server in forloop")
		if server.ready {
			var msg = <-server.ch
			server.ack = msg + 1

			client.ch <- server.ack
			client.ch <- server.seq

			<-server.ch

			client.ready = true
		}
		server.ready = false
	}
	fmt.Println("Connection established!")
	// serveren modtager beskeden fra client
	// serveren vil gerne sende sin egen seq ind i client's channel, og client's seq + 1
}
