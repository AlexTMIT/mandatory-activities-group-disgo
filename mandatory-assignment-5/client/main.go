package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	pb "replication/grpc"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var running bool
var name string
var currentAmount int

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewReplicationServiceClient(conn)

	ctx := context.Background()

	join(ctx, c)

	for {
		listenToInput(ctx, c)
	}
}

func queryBidding(ctx context.Context, c pb.ReplicationServiceClient) {
	for running {
		req, err := c.AuctionQuery(ctx, &pb.AQueryRequest{})
		if err != nil {
			log.Println("Could not fetch new messages")
		}

		log.Println(req.Result)
	}
}

func join(ctx context.Context, c pb.ReplicationServiceClient) {
	log.Println("Please input your name")
	fmt.Scanln(&name)

	req, err := c.ProcessJoinRequest(ctx, &pb.JoinRequest{ClientName: name})
	if err != nil {
		log.Println("You took too long, please try again")
	}
	log.Printf("%s", req.Msg)
	running = true
}

func bid(ctx context.Context, c pb.ReplicationServiceClient, bid int) {
	_, err := c.Bidding(ctx, &pb.BidRequest{ClientName: name, Amount: int32(bid)})
	if err != nil {
		log.Println("Error in bidding.")
	}
}

func listenToInput(ctx context.Context, c pb.ReplicationServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	command, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return
	}

	command = strings.TrimSpace(command)
	split := strings.Split(command, " ")
	keyword := strings.ToLower(split[0])

	if keyword == "query" {
		queryBidding(ctx, c)
	} else if keyword == "bid" {
		fmt.Println("How much would you like to bid?")
		fmt.Scanln(&currentAmount)
		bid(ctx, c, currentAmount)
	}
}