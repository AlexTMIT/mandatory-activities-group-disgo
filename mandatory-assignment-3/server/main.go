package main

import (
	"context"
	"fmt"
	pb "lamport_service/grpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

var lamport int32 = 0
var messages []string

type server struct {
	pb.UnimplementedChittychatServiceServer
}

func addMessage(msg string) {
	lamport++

	messages = append(messages, msg)
}

func (s *server) ProcessJoinRequest(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf("Participant %s joined Chitty-Chat at Lamport time %d", req.ParticipantName, lamport)

	msg := fmt.Sprintf("[L%d]: Welcome to ChittyChat, %s", lamport, req.ParticipantName)

	addMessage(msg)

	return &pb.JoinResponse{
		Msg:       msg,
		Timestamp: lamport,
	}, nil
}

func (s *server) ProcessLeaveRequest(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	fmt.Printf("Participant %s has left Chitty-Chat at Lamport time %d", req.ParticipantName, lamport)

	msg := fmt.Sprintf("[L%d]: See you later, %s.", lamport, req.ParticipantName)

	addMessage(msg)

	return &pb.LeaveResponse{
		Msg: msg,
	}, nil
}

func (s *server) GetMessage(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	if len(req.Msg) > 128 {
		return &pb.ChatResponse{
			Msg: "ERROR: Your message was not sent. Message was more than 128 characters.",
		}, nil
	} else {
		msg := fmt.Sprintf("%s at [L%d]: %s", req.ParticipantName, lamport, req.Msg)
		addMessage(msg)

		return &pb.ChatResponse{
			Msg: msg,
		}, nil
	}
}

func (s *server) ProcessBroadcastRequest(ctx context.Context, req *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	clientLamport := req.Timestamp
	newMessages := messages[clientLamport:]

	return &pb.BroadcastResponse{
		BroadcastMessages: newMessages,
		Timestamp:         lamport,
	}, nil
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
