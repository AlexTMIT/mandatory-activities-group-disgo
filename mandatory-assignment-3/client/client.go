package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "lamport_service/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var running bool
var name string
var lamport int32

func main() {
	conn, err := grpc.NewClient("0.0.0.0:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewChittychatServiceClient(conn)

	ctx := context.Background()

	join(ctx, c)

	go fetchNewMessages(ctx, c)

	for running {
		listenToInput(ctx, c)
	}
}

func fetchNewMessages(ctx context.Context, c pb.ChittychatServiceClient) {
	for running {
		req, err := c.ProcessBroadcastRequest(ctx, &pb.BroadcastRequest{Timestamp: lamport})
		if err != nil {
			log.Println("Could not fetch new messages")
		}
		for _, s := range req.BroadcastMessages {
			log.Println(s)
		}
		lamport = req.Timestamp
	}
}

func join(ctx context.Context, c pb.ChittychatServiceClient) {
	log.Println("Please input your name")
	fmt.Scanln(&name)

	req, err := c.ProcessJoinRequest(ctx, &pb.JoinRequest{ParticipantName: name})
	if err != nil {
		log.Println("You took too long, please try again")
	}
	log.Printf("%s", req.Msg)
	running = true
	lamport = req.Timestamp
}

func leave(ctx context.Context, c pb.ChittychatServiceClient) {
	_, err := c.ProcessLeaveRequest(ctx, &pb.LeaveRequest{ParticipantName: name})
	if err != nil {
		log.Println("Error in leaving chat.")
	}

	running = false
}

func chat(msg string, ctx context.Context, c pb.ChittychatServiceClient) {
	_, err := c.GetMessage(ctx, &pb.ChatRequest{Msg: msg, ParticipantName: name})
	if err != nil {
		log.Println("Error in sending message.")
	}
}

func listenToInput(ctx context.Context, c pb.ChittychatServiceClient) {
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
