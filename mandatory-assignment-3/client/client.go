package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
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
		log.Println("You took too long, please try again")
	}
	log.Printf("%s\n", req.Msg)
	running = true
}

func leave(ctx context.Context, c pb.ChittychatServiceClient) {
	req, err := c.ProcessLeaveRequest(ctx, &pb.LeaveRequest{ParticipantName: name})
	if err != nil {
		log.Println("You took too long, please try again.")
	}

	log.Printf("%s\n", req.Msg)
	running = false
}

func chat(msg string, ctx context.Context, c pb.ChittychatServiceClient) {
	req, err := c.GetMessage(ctx, &pb.ChatRequest{Msg: msg, ParticipantName: name})
	if err != nil {
		log.Println("Error in sending message.")
	}
	log.Printf("%s\n", req.Msg)
}

func listen(ctx context.Context, c pb.ChittychatServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	command, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return
	}

	command = strings.TrimSpace(command)
	split := strings.Split(command, " ")
	keyword := strings.ToLower(split[0])

	if keyword == "leave" {
		leave(ctx, c)
	} else if keyword == "chat" {
		msg := strings.Join(split[1:], " ")
		chat(msg, ctx, c)
	}
}
