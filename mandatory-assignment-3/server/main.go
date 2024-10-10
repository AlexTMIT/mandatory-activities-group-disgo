package main

import (
	"context"
	"fmt"
	pb "lamport_service/grpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChittychatServiceServer
}

func (s *server) ProcessJoinRequest(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf("Participant %s joined Chitty-Chat at Lamport time L\n", req.ParticipantName)

	return &pb.JoinResponse{
		Msg: fmt.Sprintf("Welcome to ChittyChat, %s\n", req.ParticipantName),
	}, nil
}

func (s *server) ProcessLeaveRequest(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	fmt.Printf("Participant %s has left Chitty-Chat at Lamport time L\n", req.ParticipantName)

	return &pb.LeaveResponse{
		Msg: fmt.Sprintf("See you later, %s\n", req.ParticipantName),
	}, nil
}

func (s *server) GetMessage(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	if len(req.Msg) > 128 {
		return &pb.ChatResponse{
			Msg: "ERROR: Your message was not sent. Reason: message was longer than 128 characters.\n",
		}, nil
	} else {
		return &pb.ChatResponse{
			Msg: fmt.Sprintf("%s: %s\n", req.ParticipantName, req.Msg),
		}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChittychatServiceServer(s, &server{})
	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
