package process

import (
	pb "consensus/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type vars struct {
	isInsideCS   bool
	isCSEmpty    bool
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
}

type process struct {
	pb.UnimplementedConsensusServiceServer
	vars *vars
}

func newProcess() *process {
	return &process{
		vars: &vars{
			isInsideCS:   false,
			isCSEmpty:    true,
			currentState: RELEASED,
		},
	}
}

func (s *process) CriticalSection(ctx context.Context, req *pb.CriticalRequest) (*pb.CriticalReply, error) {
	// update lamport
	if req.Lamport > s.vars.lamport {
		s.vars.lamport = req.Lamport
	}
	s.vars.lamport++

	// decide whether to grant access
	grant := false

	if s.vars.currentState == RELEASED {
		grant = true
	} else if s.vars.currentState == HELD {
		grant = false
	} else if s.vars.currentState == WANTED {
		if req.Lamport < s.vars.lamport || (req.Lamport == s.vars.lamport && req.Port < s.vars.id) {
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

func (s *process) JoiningQueue(ctx context.Context, req *pb.JoiningRequest) (*pb.JoiningReply, error) {
	s.vars.requests = append(s.vars.requests, queueItem{id: req.Port, lamport: req.Lamport})
	s.SortClientsByLamport()

	return &pb.JoiningReply{}, nil
}

func (s *process) EnteringCS(ctx context.Context, req *pb.EnteringCSRequest) (*pb.EnteringCSReply, error) {
	s.vars.isCSEmpty = false

	return &pb.EnteringCSReply{}, nil
}

func (s *process) ExitingCS(ctx context.Context, req *pb.ExitingCSRequest) (*pb.ExitingCSReply, error) {
	s.vars.isCSEmpty = true
	s.vars.requests = s.vars.requests[1:]

	return &pb.ExitingCSReply{}, nil
}

func (s *process) GetLamport(ctx context.Context, req *pb.LamportRequest) (*pb.LamportReply, error) {
	return &pb.LamportReply{Lamport: s.vars.lamport}, nil
}

func (s *process) SortClientsByLamport() {
	sort.Slice(s.vars.clients, func(i, j int) bool {
		return s.vars.requests[i].lamport < s.vars.requests[j].lamport
	})
}

func (s *process) inRequest() bool {
	for _, e := range s.vars.requests {
		if e.id == s.vars.id {
			return true
		}
	}
	return false
}

func (s *process) createClients() {
	for i := 0; i < len(s.vars.ports); i++ {
		ctx, c := createClient(s.vars.ports[i])
		s.vars.clients = append(s.vars.clients, client{ctx: ctx, c: c})
	}
}

func (s *process) multicastCSRequest() {
	for i := 0; i < len(s.vars.clients); i++ {
		var client = s.vars.clients[i]
		s.makeRequest(client.ctx, client.c)
	}
}

func createClient(port string) (ctx context.Context, c pb.ConsensusServiceClient) {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	c = pb.NewConsensusServiceClient(conn)

	ctx = context.Background()
	return
}

func (s *process) makeRequest(ctx context.Context, c pb.ConsensusServiceClient) {
	_, err := c.CriticalSection(ctx, &pb.CriticalRequest{Port: s.vars.id, Lamport: s.vars.lamport})
	if err != nil {
		log.Printf("Error in making request: %v\n", err)
	}
}

func (s *process) initProcessServer() {
	lis, err := net.Listen("tcp", s.vars.serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterConsensusServiceServer(grpcServer, s)
	log.Printf("Server is running on port %s...\n", s.vars.serverPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run process: %v", err)
	}
}

func (s *process) checkReplies() {
	if s.vars.replies == int32(len(s.vars.ports)-1) {
		s.vars.currentState = HELD
		for i := 0; i < len(s.vars.clients); i++ {
			s.multicastJoiningRequest(i)
		}
	}
}

func (s *process) multicastEnteringRequest(i int) {
	_, err := s.vars.clients[i].c.EnteringCS(s.vars.clients[i].ctx, &pb.EnteringCSRequest{})
	if err != nil {
		log.Printf("Error in multicastEnteringRequest: %v\n", err)
	}
}

func (s *process) multicastExitingRequest(i int) {
	_, err := s.vars.clients[i].c.ExitingCS(s.vars.clients[i].ctx, &pb.ExitingCSRequest{})
	if err != nil {
		log.Println("Error in multicastExitingRequest:", err)
	}
}

func (s *process) multicastLamportRequest(i int) {
	rep, err := s.vars.clients[i].c.GetLamport(s.vars.clients[i].ctx, &pb.LamportRequest{})
	if err != nil {
		log.Println("Error in multicastLamportRequest:", err)
		return
	}

	if rep.Lamport > s.vars.lamport {
		s.vars.lamport = rep.Lamport + 1
	}
}

func (s *process) multicastJoiningRequest(i int) {
	_, err := s.vars.clients[i].c.JoiningQueue(s.vars.clients[i].ctx, &pb.JoiningRequest{Port: s.vars.id, Lamport: s.vars.lamport})
	if err != nil {
		log.Println("Error in multicastJoiningRequest:", err)
	}
}

func (s *process) goInsideCS() {
	if len(s.vars.requests) == 0 {
		return
	}

	if s.vars.requests[0].id == s.vars.id && s.vars.isCSEmpty {
		for i := 0; i < len(s.vars.clients); i++ {
			s.multicastEnteringRequest(i)
			s.multicastLamportRequest(i)
			fmt.Printf("Client %d has entered CS", s.vars.id)
			s.multicastExitingRequest(i)
			fmt.Printf("Client %d has left CS", s.vars.id)
		}
	}
}

func (s *process) shouldMulticast() bool {
	if s.vars.isInsideCS {
		return false
	}
	if s.inRequest() {
		return false
	}
	if s.vars.replies > 0 && int(s.vars.replies) < len(s.vars.ports)-1 {
		return false
	}
	return true
}

func (s *process) initialize(porto string, portList []string) {
	s.vars.serverPort = porto
	s.vars.ports = portList
	s.vars.currentState = RELEASED

	p, _ := strconv.Atoi(porto[len("localhost:"):])

	fmt.Printf("parsed id %d from serverPort %s\n", p, porto)

	s.vars.id = int32(p)

	s.createClients()
}

func Run(porto string, portList []string) {
	s := newProcess()
	s.initialize(porto, portList)

	// start the server
	go s.initProcessServer()

	// start other goroutines
	go s.checkReplies()
	go s.goInsideCS()

	// main loop
	for {
		if s.shouldMulticast() {
			s.vars.replies = 0
			s.multicastCSRequest()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
