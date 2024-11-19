package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "replication/grpc"

	"google.golang.org/grpc"
)

var bidAmount int32 = 0
var finished bool
var highestBidder string
var amountOfBids int = 0

type server struct {
	pb.UnimplementedReplicationServiceServer
}

func (s *server) ProcessJoinRequest(ctx context.Context, req *pb.JoinRequest) (*pb.JoinReply, error) {
	msg := fmt.Sprintf("Welcome to the auction, %s!", req.ClientName)

	return &pb.JoinReply{
		Msg: msg,
	}, nil
}

func (s *server) Bidding(ctx context.Context, req *pb.BidRequest) (*pb.BidReply, error) {
	var msg = "Hello"
	if amountOfBids == 10 || finished {
		msg = "The bidding has ended."
		finished = true
	}
	if req.Amount > bidAmount && !finished {
		highestBidder = req.ClientName
		msg = fmt.Sprintf("Client %s is bidding with an amount of %d", req.ClientName, req.Amount)
		bidAmount = req.Amount
		amountOfBids++
	} else if !finished {
		msg = "Please enter an amount higher than the current bidding."
	}

	fmt.Println(msg)
	return &pb.BidReply{
		Response: msg,
	}, nil
}

func (s *server) AuctionQuery(ctx context.Context, req *pb.AQueryRequest) (*pb.AQueryReply, error) {
	var msg = ""
	if finished {
		msg = fmt.Sprintf("The bidding has finished on a total amount of %d, with %s as the winner!", bidAmount, highestBidder)
		fmt.Println(msg)
	} else {
		msg = fmt.Sprintf("Current Bidding Amount is on %d", bidAmount)
	}

	return &pb.AQueryReply{
		CurrentAmount: bidAmount,
		Result:        msg,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterReplicationServiceServer(s, &server{})
	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
