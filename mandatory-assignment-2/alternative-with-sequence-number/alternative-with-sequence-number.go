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
	fmt.Printf("running client. client seq starts at %d\n", client.seq)
	for !finished {
		if !client.ready {
			server.ch <- client.seq
			server.ready = true
		} else {
			client.seq = <-client.ch //server's ack no.
			//<-client.ch

			var msg = <-client.ch
			client.ack = msg + 1

			server.ch <- client.seq
			server.ch <- client.ack
			fmt.Println("Almost finished")
			finished = true
			break
		}
	}
	fmt.Printf("Client connection established! Client seq is %d. Client ack is %d\n", client.seq, client.ack)
	// client vil sende seq nummer ind til server
	// client skal herefter vente til, at client får 2 numre ind i sin egen kanal
	// når client får besken, vil den gerne sende serveren's seq + 1, plus sin egen seq nummer + 1
}

func serverGo(client *client, server *server) {
	fmt.Printf("server seq starts at %d\n", server.seq)
	for !finished {
		if server.ready {
			client.ready = true
			var msg = <-server.ch
			server.ack = msg + 1

			client.ch <- server.ack
			client.ch <- server.seq

			<-server.ch

			server.ready = false
		}
	}
	fmt.Printf("Server connection established! Server seq is %d. Server ack is %d\n", server.seq, server.ack)

}
