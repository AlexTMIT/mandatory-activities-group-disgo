package process

import (
	pb "consensus/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var id int32
var serverPort string
var lamport int32
var ports []string
var currentState State
var requests []queueItem
var replies int32
var clients []client

type queueItem struct {
	port    int32
	lamport int32
}

type client struct {
	ctx context.Context
	c   pb.ConsensusServiceClient
}

type process struct {
	pb.UnimplementedConsensusServiceServer
}

func (s *process) ProcessConsensus(ctx context.Context, req *pb.CriticalRequest) (*pb.CriticalReply, error) {
	currentState = WANTED
	requests = append(requests, queueItem{port: req.Port, lamport: req.Lamport})
	replies++

	fmt.Printf("Process %d is requesting to join Critical Section at Lamport time %d", req.Port, req.Lamport)

	return &pb.CriticalReply{}, nil
}

func Run(porto string, portList []string) {
	initialize(porto, portList)

	for {
		if !inRequest() {
			replies = 0
			broadcastCSRequest()
		}
	}
}

func inRequest() bool {
	for _, e := range requests {
		if e.port == id {
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

func broadcastCSRequest() {
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
