package main

import (
	"context"
	"fmt"
	pb "lamport_service/grpc"
	//"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChittychatServiceServer
}

func (s *server) ProcessJoinRequest(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf("Participant %s joined Chitty-Chat at Lamport time L", req.ParticipantName)

	return &pb.JoinResponse{
		Msg: "Successfully joined chittychat!",
	}, nil
}

func (s *server) ProcessLeaveRequest(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	fmt.Printf("Participant %s has left Chitty-Chat at Lamport time L", req.ParticipantName)

	return &pb.LeaveResponse{
		Msg: "Successfully left chittychat!",
	}, nil
}

func (s *server) GetMessage(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	return &pb.ChatResponse{
		Msg: fmt.Sprintf("Participant joined Chitty-Chat at Lamport time L"),
	}, nil
}
