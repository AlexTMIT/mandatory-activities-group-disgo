package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	pb "replication/grpc"

	"google.golang.org/grpc"
)

var lamport int = 0
var bidAmount int32 = 0
var finished bool
var highestBidder string
var port string

type server struct {
	pb.UnimplementedReplicationServiceServer
}

func (s *server) Bidding(ctx context.Context, req *pb.BidRequest) (*pb.BidReply, error) {
	lamport++
	if lamport < int(req.Lamport) {
		log.Println("This server's lamport was lower than the bidder")
		log.Println("The server has previously crashed")

		return nil, errors.New("server is lacking behind, and will therefore not respond")
	}

	var msg string
	log.Printf("Lamport %d server %s", lamport, port)
	if lamport == 7 || finished {
		msg = "FAIL: The bidding has ended."
		fmt.Println(msg[6:])
		finished = true
	}
	if req.Amount > bidAmount && !finished {
		highestBidder = req.ClientName

		msg = "SUCCESS"
		fmt.Printf("Client %s is bidding with an amount of %d \n", req.ClientName, req.Amount)
		bidAmount = req.Amount
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
	log.Println("0.0.0.0:500 + ?: ")
	fmt.Scanln(&port)

	lis, err := net.Listen("tcp", "0.0.0.0:500"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterReplicationServiceServer(s, &server{})
	log.Println("Server is running on port 500" + port + "...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
