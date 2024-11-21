package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	pb "replication/grpc"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name string
var ctx1 context.Context
var c1 pb.ReplicationServiceClient
var ctx2 context.Context
var c2 pb.ReplicationServiceClient
var responses []string

func main() {
	setName()
	processJoin()

	for {
		listenToInput()
	}
}

func processJoin() {
	ctx1, c1 = connectToServer("localhost:50051")
	ctx2, c2 = connectToServer("localhost:50052")

	time.Sleep(time.Second)
	printResponse()
}

func connectToServer(port string) (ctx context.Context, c pb.ReplicationServiceClient) {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	c = pb.NewReplicationServiceClient(conn)

	ctx = context.Background()
	join(ctx, c)

	return
}

func setName() {
	log.Println("Please input your name")
	fmt.Scanln(&name)
}

func queryBidding(ctx context.Context, c pb.ReplicationServiceClient) {
	req, err := c.AuctionQuery(ctx, &pb.AQueryRequest{})
	if err != nil {
		log.Println("Could not query")
	}
	log.Println(req.Result)
}

func join(ctx context.Context, c pb.ReplicationServiceClient) {
	req, err := c.ProcessJoinRequest(ctx, &pb.JoinRequest{ClientName: name})
	if err != nil {
		log.Println("Failed to process join request")
	}
	responses = append(responses, req.Msg)
}

func bid(ctx context.Context, c pb.ReplicationServiceClient, bid int) {
	req, err := c.Bidding(ctx, &pb.BidRequest{Amount: int32(bid), ClientName: name})
	if err != nil {
		log.Println("One server is down.")
		return
	}
	responses = append(responses, req.Response)
}

func listenToInput() {
	command := readCommand()
	split, keyword := parseCommand(command)

	if keyword == "query" {
		queryBidding(ctx1, c1)
		queryBidding(ctx2, c2)
	} else if keyword == "bid" {
		processBid(split)
	}
}

func processBid(split []string) {
	currentAmount, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("Invalid input.")
	}

	bid(ctx1, c1, currentAmount)
	bid(ctx2, c2, currentAmount)

	time.Sleep(time.Second)
	printResponse()
}

func printResponse() {
	if len(responses) > 0 {
		log.Println(responses[0])
		responses = responses[:0]
	} else {
		log.Println("hey... all servers are down :(")
	}
}

func parseCommand(command string) ([]string, string) {
	split := strings.Split(command, " ")
	keyword := strings.ToLower(split[0])
	return split, keyword
}

func readCommand() string {
	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return ""
	}

	return strings.TrimSpace(command)
}
