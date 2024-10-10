package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pb "lamport_service/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var running bool
var name string

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewChittychatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	join(ctx, c)
	for running {
		listen(ctx, c)
		time.Sleep(1 * time.Second)
	}
}

func join(ctx context.Context, c pb.ChittychatServiceClient) {
	log.Println("Please input your name")
	fmt.Scanln(&name)

	req, err := c.ProcessJoinRequest(ctx, &pb.JoinRequest{ParticipantName: name})
	if err != nil {
		log.Print("You took too long, please try again")
	}
	log.Printf(req.Msg)
	running = true
}

func leave(ctx context.Context, c pb.ChittychatServiceClient) {
	req, err := c.ProcessLeaveRequest(ctx, &pb.LeaveRequest{ParticipantName: name})
	if err != nil {
		log.Print("You took too long, please try again")
	}

	log.Printf(req.Msg)
	running = false
}

func listen(ctx context.Context, c pb.ChittychatServiceClient) {
	var command string
	fmt.Scanln(&command)
	split := strings.Split(command, " ")

	if strings.ToLower(split[0]) == "leave" {
		leave(ctx, c)
	}
}
