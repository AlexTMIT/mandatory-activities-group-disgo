package process

import (
	pb "consensus/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type vars struct {
	mu           sync.Mutex
	isInsideCS   bool
	id           int32
	serverPort   string
	lamport      int32
	ports        []string
	currentState State
	requests     []queueItem
	replies      int32
	clients      []client
}

type queueItem struct {
	id      int32
	lamport int32
}

type client struct {
	ctx context.Context
	c   pb.ConsensusServiceClient
	id  int32
}

type process struct {
	pb.UnimplementedConsensusServiceServer
	vars *vars
}

func newProcess() *process {
	return &process{
		vars: &vars{
			isInsideCS:   false,
			currentState: RELEASED,
		},
	}
}

func (s *process) CriticalSection(ctx context.Context, req *pb.CriticalRequest) (*pb.CriticalReply, error) {
	s.vars.mu.Lock()
	defer s.vars.mu.Unlock()

	// update lamport clock
	if req.Lamport > s.vars.lamport {
		s.vars.lamport = req.Lamport
	}
	s.vars.lamport++

	// decide whether to grant access
	var grant bool
	if s.vars.currentState == RELEASED {
		grant = true
	} else if s.vars.currentState == HELD {
		grant = false
	} else if s.vars.currentState == WANTED {
		if req.Lamport > s.vars.lamport || req.Port > s.vars.id {
			grant = false
		} else {
			grant = true
		}
	}

	if grant {
		return &pb.CriticalReply{Grant: true}, nil

	} else {
		// queue the request
		s.vars.requests = append(s.vars.requests, queueItem{id: req.Port, lamport: req.Lamport})
		return &pb.CriticalReply{Grant: false}, nil
	}
}

func (s *process) ReplyCS(ctx context.Context, req *pb.ReplyRequest) (*pb.ReplyReply, error) {
	s.vars.mu.Lock()
	s.vars.replies++
	s.vars.mu.Unlock()
	return &pb.ReplyReply{}, nil
}

func (s *process) sendDeferredReplies() {
	deferredRequests := s.vars.requests
	s.vars.requests = nil

	for _, req := range deferredRequests {
		s.sendReply(req)
	}
}
func (s *process) sendReply(req queueItem) {
	// Find the client for the requesting process
	var client *client
	for i := range s.vars.clients {
		c := &s.vars.clients[i]
		if c.id == req.id {
			client = c
			break
		}
	}

	if client != nil {
		_, err := client.c.ReplyCS(client.ctx, &pb.ReplyRequest{Port: s.vars.id})
		if err != nil {
			log.Printf("error sending deferred reply to %d: %v\n", req.id, err)
		}
	} else {
		log.Printf("client not found for process %d", req.id)
	}
}

func (s *process) createClients() {
	for i := 0; i < len(s.vars.ports); i++ {
		portStr := s.vars.ports[i]
		if portStr != s.vars.serverPort {
			ctx, c := createClient(portStr)

			// Extract the process ID from the port string
			p, _ := strconv.Atoi(portStr[len("localhost:"):])
			clientID := int32(p)

			s.vars.clients = append(s.vars.clients, client{ctx: ctx, c: c, id: clientID})
		}
	}
}

func createClient(port string) (ctx context.Context, c pb.ConsensusServiceClient) {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c = pb.NewConsensusServiceClient(conn)
	ctx = context.Background()
	return
}

func (s *process) makeRequest(ctx context.Context, c pb.ConsensusServiceClient) {
	resp, err := c.CriticalSection(ctx, &pb.CriticalRequest{Port: s.vars.id, Lamport: s.vars.lamport})
	if err != nil {
		log.Printf("error in making request: %v\n", err)
	} else if resp.Grant {
		s.vars.mu.Lock()
		s.vars.replies++
		s.vars.mu.Unlock()
	}
}

func (s *process) multicastCSRequest() {
	for i := 0; i < len(s.vars.clients); i++ {
		var client = s.vars.clients[i]
		s.makeRequest(client.ctx, client.c)
	}
}

func (s *process) initProcessServer() {
	lis, err := net.Listen("tcp", s.vars.serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterConsensusServiceServer(grpcServer, s)
	log.Printf("server is running on port %s...\n", s.vars.serverPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run process: %v", err)
	}
}

func (s *process) checkReplies() {
	for {
		s.vars.mu.Lock()
		if s.vars.currentState == WANTED && s.vars.replies == int32(len(s.vars.clients)-1) {
			s.vars.currentState = HELD
			s.vars.isInsideCS = true
			s.vars.mu.Unlock()

			// Simulate entering the critical section
			fmt.Printf("process %d has entered CS\n", s.vars.id)
			fmt.Printf("process %d has left CS\n", s.vars.id)

			s.vars.mu.Lock()
			s.vars.isInsideCS = false
			s.vars.currentState = RELEASED
			s.vars.replies = 0
			s.vars.mu.Unlock()

			// Send any deferred replies to other processes
			s.sendDeferredReplies()
		} else {
			s.vars.mu.Unlock()
		}
		//time.Sleep(100 * time.Millisecond)
	}
}

func (s *process) shouldMulticast() bool {
	s.vars.mu.Lock()
	defer s.vars.mu.Unlock()

	return s.vars.currentState == RELEASED
}

func (s *process) initialize(porto string, portList []string) {
	s.vars.serverPort = porto
	s.vars.ports = portList
	s.vars.currentState = RELEASED

	p, _ := strconv.Atoi(porto[len("localhost:"):])
	s.vars.id = int32(p)

	// start the server
	go s.initProcessServer()

	// wait to ensure the server is running
	time.Sleep(2 * time.Second)

	s.createClients()
}

func Run(porto string, portList []string) {
	s := newProcess()
	s.initialize(porto, portList)

	// start other goroutines
	go s.checkReplies()

	// main loop
	for {
		if s.shouldMulticast() {
			s.vars.mu.Lock()
			s.vars.currentState = WANTED
			s.vars.replies = 0
			go s.multicastCSRequest()
			s.vars.lamport++
			s.vars.mu.Unlock()

		}
		time.Sleep(100 * time.Millisecond)
	}
}
