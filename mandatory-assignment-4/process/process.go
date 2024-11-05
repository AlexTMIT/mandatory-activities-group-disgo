package process

import (
	pb "consensus/grpc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var port string
var id int32
var ports []string

type process struct {
	pb.UnimplementedConsensusServiceServer
}

func (s *process) ProcessConsensus(ctx context.Context, req *pb.CriticalRequest) (*pb.CriticalReply, error) {
	fmt.Printf("Process %d is requesting to join Critical Section at Lamport time %d", req.Id, req.Lamport)

	return &pb.CriticalReply{}, nil
}

func Run(porto string, idi int32, portList []string) {
	initialize(porto, idi, portList)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConsensusServiceServer(s, &process{})
	log.Printf("Server is running on port %s...\n", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to run process: %v", err)
	}
}

func initialize(porto string, idi int32, portList []string) {
	port = porto
	id = idi
	ports = portList
}
