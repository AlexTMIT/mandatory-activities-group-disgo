package main

import (
	"fmt"
	"math/rand"
	"time"
)

type msgType int
type stateType int

const (
	SYN msgType = iota
	SYN_ACK
	ACK
)

const (
	CLOSED stateType = iota
	LISTEN
	SYN_SENT
	SYN_RECEIVED
	ESTABLISHED
)

type message struct {
	msgType  msgType
	clientID int
	serverID int
}

type client struct {
	id             int
	state          stateType
	clientToServer chan message
	serverToClient chan message
}

type server struct {
	id             int
	state          stateType
	clientChannels map[int]chan message // for client to server
	serverChannels map[int]chan message // for server to client
}

var clientDone []bool

// client logic ------------------------------------------------------------------------------

func clientGo(c *client) {
	fmt.Printf("Starting client %d\n", c.id)
	sendSYN(c)
	waitSYN_ACK(c)
}

func sendSYN(c *client) {
	c.clientToServer <- message{msgType: SYN, clientID: c.id}
	c.state = SYN_SENT
	fmt.Printf("Client %d sent SYN\n", c.id)
}

func waitSYN_ACK(c *client) {
	for c.state != ESTABLISHED {
		select {
		case msg := <-c.serverToClient:
			receiveSYN_ACK(msg, c)
			sendACK(c, msg.serverID)
			clientDone[c.id] = true
		case <-time.After(5 * time.Second):
			timeoutSYN_ACK(c)
		}
	}
}

func sendACK(c *client, serverID int) {
	c.clientToServer <- message{msgType: ACK, clientID: c.id}
	fmt.Printf("Client %d sent ACK to server %d\n", c.id, serverID)
}

func receiveSYN_ACK(msg message, c *client) {
	if msg.msgType != SYN_ACK {
		return
	}

	fmt.Printf("Client %d received SYN-ACK from server %d\n", c.id, msg.serverID)
	c.state = ESTABLISHED
}

func timeoutSYN_ACK(c *client) {
	fmt.Printf("Client %d had timeout waiting for SYN-ACK\n", c.id)
	sendSYN(c)
}

// server logic ------------------------------------------------------------------------------

func serverGo(s *server) {
	fmt.Printf("Starting server %d\n", s.id)
	s.state = LISTEN

	for {
		runServer(s)
	}
}

func runServer(s *server) {
	for clientID, clientChannel := range s.clientChannels {
		select {
		case msg := <-clientChannel:
			if msg.msgType == SYN {
				receiveSYN(s, msg)
				sendSYN_ACK(s, clientID)
				waitACK(s, clientID, clientChannel)
			}
		}
	}
}

func waitACK(s *server, clientID int, clientChannel chan message) {
	select {
	case msg := <-clientChannel:
		if msg.msgType == ACK {
			fmt.Printf("Server %d received ACK from Client %d, connection established\n", s.id, clientID)
			s.state = ESTABLISHED
		}
	case <-time.After(5 * time.Second):
		fmt.Printf("Server %d: Timeout waiting for ACK from Client %d, resending SYN-ACK\n", s.id, clientID)
		sendSYN_ACK(s, clientID)
	}
}

func sendSYN_ACK(s *server, clientID int) {
	clientChannel := s.serverChannels[clientID]
	clientChannel <- message{msgType: SYN_ACK, clientID: clientID, serverID: s.id}
	fmt.Printf("Server %d sent SYN-ACK to client %d\n", s.id, clientID)
}

func receiveSYN(s *server, msg message) {
	fmt.Printf("Server %d received SYN from client %d\n", s.id, msg.clientID)
	s.state = SYN_RECEIVED
}

// main logic ------------------------------------------------------------------------------

func initClients(NUM_CLIENTS int, NUM_SERVERS int, servers []*server) {
	for i := 0; i < NUM_CLIENTS; i++ {
		c := initClient(NUM_SERVERS, servers, i)
		go clientGo(c)
	}
}

func initClient(NUM_SERVERS int, servers []*server, i int) *client {
	clientToServer := make(chan message)
	serverToClient := make(chan message)

	serverID := rand.Intn(NUM_SERVERS) // assign random server
	server := servers[serverID]

	server.clientChannels[i] = clientToServer
	server.serverChannels[i] = serverToClient

	return &client{
		id:             i,
		state:          CLOSED,
		clientToServer: clientToServer,
		serverToClient: serverToClient,
	}
}

func initServers(NUM_SERVERS int, servers []*server) {
	for i := 0; i < NUM_SERVERS; i++ {
		servers[i] = &server{
			id:             i,
			state:          CLOSED,
			clientChannels: make(map[int]chan message),
			serverChannels: make(map[int]chan message),
		}

		go serverGo(servers[i])
	}
}

func allClientsDone() bool {
	for _, done := range clientDone {
		if !done {
			return false
		}
	}
	return true
}

func main() {
	NUM_SERVERS := 10
	NUM_CLIENTS := 50

	servers := make([]*server, NUM_SERVERS)
	clientDone = make([]bool, NUM_CLIENTS)

	initServers(NUM_SERVERS, servers)
	initClients(NUM_CLIENTS, NUM_SERVERS, servers)

	for {
		if allClientsDone() {
			fmt.Println("All clients are done!")
			break
		}
		time.Sleep(1 * time.Second)
	}
}
