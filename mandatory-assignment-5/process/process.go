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
	msg := fmt.Sprintf("Welcome to the auction, %s! \nTo bid an amount, type 'bid'.\nTo find the current highest bidder, type 'query'.", req.ClientName)

	return &pb.JoinReply{
		Msg: msg,
	}, nil
}

func (s *server) Bidding(ctx context.Context, req *pb.BidRequest) (*pb.BidReply, error) {
	var msg = "Hello"
	if amountOfBids == 10 || finished {
		msg = fmt.Sprintln("FAIL:" + "The bidding has ended.")
		finished = true
	}
	if req.Amount > bidAmount && !finished {
		highestBidder = req.ClientName

		msg = "SUCCESS"
		fmt.Printf("Client %s is bidding with an amount of %d \n", req.ClientName, req.Amount)
		bidAmount = req.Amount
		amountOfBids++
	} else if !finished {
		msg = "FAIL: Please enter an amount higher than the current bidding."
		fmt.Printf("Client %s has entered a bidding amount too low.\n", req.ClientName)
	}

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
