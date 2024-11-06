package process

import (
	pb "consensus/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"

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

func (s *process) ProcessConsensus(ctx context.Context, req *pb.CriticalRequest) (*pb.CriticalReply, error) {
	currentState = WANTED
	replies++

	fmt.Printf("Process %d is requesting to join Critical Section at Lamport time %d", req.Port, req.Lamport)

	return &pb.CriticalReply{}, nil
}

func (s *process) JoiningQueue(ctx context.Context, req *pb.JoiningRequest) (*pb.JoiningReply, error) {
	requests = append(requests, queueItem{id: req.Port, lamport: req.Lamport})
	SortClientsByLamport()

	return &pb.JoiningReply{}, nil
}

func (s *process) EnteringCS(ctx context.Context, req *pb.EnteringCSRequest) (*pb.EnteringCSReply, error) {
	isCSEmpty = false

	return &pb.EnteringCSReply{}, nil
}

func (s *process) ExitingCS(ctx context.Context, req *pb.ExitingCSRequest) (*pb.ExitingCSReply, error) {
	isCSEmpty = true
	requests = requests[1:]

	return &pb.ExitingCSReply{}, nil
}

func SortClientsByLamport() {
	sort.Slice(clients, func(i, j int) bool {
		return requests[i].lamport < requests[j].lamport
	})
}

func inRequest() bool {
	for _, e := range requests {
		if e.id == id {
			return true
		}
	}

	return false
}

func createClients() {
	for i := 0; i < len(ports); i++ {
		ctx, c := createClient(ports[i])
		clients = append(clients, client{ctx: ctx, c: c})
	}
}

func multicastCSRequest() {
	for i := 0; i < len(clients); i++ {
		var client = clients[i]
		makeRequest(client.ctx, client.c)
	}
}

func createClient(port string) (ctx context.Context, c pb.ConsensusServiceClient) {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c = pb.NewConsensusServiceClient(conn)

	ctx = context.Background()
	return
}

func makeRequest(ctx context.Context, c pb.ConsensusServiceClient) {
	_, err := c.CriticalSection(ctx, &pb.CriticalRequest{Port: id, Lamport: lamport})
	if err != nil {
		log.Println("You took too long, please try again")
	}
}

func initProcessServer() {
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConsensusServiceServer(s, &process{})
	log.Printf("Server is running on port %s...\n", serverPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to run process: %v", err)
	}
}

func initialize(porto string, portList []string) {
	serverPort = porto
	ports = portList
	currentState = RELEASED

	p, _ := strconv.Atoi(serverPort[:5])
	id = int32(p)

	createClients()
	initProcessServer()
}

func checkReplies() {
	if replies == int32(len(ports)-1) {
		currentState = HELD
		for i := 0; i < len(clients); i++ {
			multicastJoiningRequest(i)
		}
	}
}

func multicastEnteringRequest(i int) {
	_, err := clients[i].c.EnteringCS(clients[i].ctx, &pb.EnteringCSRequest{})
	if err != nil {
		log.Println("You took too long, please try again")
	}
}

func multicastExitingRequest(i int) {
	_, err := clients[i].c.ExitingCS(clients[i].ctx, &pb.ExitingCSRequest{})
	if err != nil {
		log.Println("You took too long, please try again")
	}
}

func multicastLamportRequest(i int) {
	rep, err := clients[i].c.GetLamport(clients[i].ctx, &pb.LamportRequest{})
	if err != nil {
		log.Println("You took too long, please try again")
	}

	if rep.Lamport > lamport {
		lamport = rep.Lamport + 1
	}
}

func multicastJoiningRequest(i int) {
	_, err := clients[i].c.JoiningQueue(clients[i].ctx, &pb.JoiningRequest{Port: id, Lamport: lamport})
	if err != nil {
		log.Println("You took too long, please try again")
	}
}

func insideCS() {
	if requests[0].id == int32(id) && isCSEmpty {
		for i := 0; i < len(clients); i++ {
			multicastEnteringRequest(i)
			multicastLamportRequest(i)
			fmt.Printf("Client %d has entered CS", id)
			multicastExitingRequest(i)
			fmt.Printf("Client %d has left CS", id)
		}
	}
}

func shouldMulticast() bool {
	if isInsideCS {
		return false
	}

	if inRequest() {
		return false
	}

	if replies > 0 && int(replies) < len(ports)-1 {
		return false
	}

	return true
}

func Run(porto string, portList []string) {
	initialize(porto, portList)

	go checkReplies()
	go insideCS()

	for {
		if shouldMulticast() {
			replies = 0
			multicastCSRequest()
		}
	}
}
