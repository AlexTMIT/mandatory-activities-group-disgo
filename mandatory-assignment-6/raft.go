package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"time"
)

type node struct {
	connect bool
	address string
}

func newNode(address string) *node {
	node := &node{}
	node.address = address
	return node
}

// State def
type State int

// status of node
const (
	Follower State = iota + 1
	Candidate
	Leader
)

// LogEntry struct
type LogEntry struct {
	LogTerm  int
	LogIndex int
	LogCMD   interface{}
}

// Raft Node
type Raft struct {
	id int

	peers map[int]*node

	state       State
	currentTerm int
	votedFor    int
	voteCount   int

	log []LogEntry

	commitIndex int
	lastApplied int

	nextIndex  []int
	matchIndex []int

	heartbeatC chan bool
	leaderC    chan bool
}

// RequestVote rpc method
func (rf *Raft) RequestVote(args VoteArgs, reply *VoteReply) error {

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return nil
	}

	if rf.votedFor == -1 {
		rf.currentTerm = args.Term
		rf.votedFor = args.CandidateID
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	}

	return nil
}

// Heartbeat rpc method
func (rf *Raft) Heartbeat(args HeartbeatArgs, reply *HeartbeatReply) error {

	if args.Term < rf.currentTerm {
		reply.Success = false
		reply.Term = rf.currentTerm
		return nil
	}

	rf.heartbeatC <- true
	if len(args.Entries) == 0 {
		reply.Success = true
		reply.Term = rf.currentTerm
		return nil
	}

	if args.PrevLogIndex > rf.getLastIndex() {
		reply.Success = false
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
		return nil
	}

	rf.log = append(rf.log, args.Entries...)
	rf.commitIndex = rf.getLastIndex()
	reply.Success = true
	reply.Term = rf.currentTerm
	reply.NextIndex = rf.getLastIndex() + 1

	return nil
}

func (rf *Raft) rpc(port string) {
	rpc.Register(rf)
	rpc.HandleHTTP()
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal("listen error: ", err)
		}
	}()
}

func (rf *Raft) start() {
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartbeatC = make(chan bool)
	rf.leaderC = make(chan bool)

	go func() {

		rand.Seed(time.Now().UnixNano())

		for {
			switch rf.state {
			case Follower:
				select {
				case <-rf.heartbeatC:
					log.Printf("follower-%d recived heartbeat\n", rf.id)
				case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond): // sussy time interval
					log.Printf("follower-%d timeout\n", rf.id)
					rf.state = Candidate
				}
			case Candidate:
				fmt.Printf("Node: %d, I'm candidate\n", rf.id)
				rf.currentTerm++
				rf.votedFor = rf.id
				rf.voteCount = 1
				go rf.broadcastRequestVote()

				select {
				case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
					rf.state = Follower
				case <-rf.leaderC:
					fmt.Printf("Node: %d, I'm leader\n", rf.id)
					rf.state = Leader

					rf.nextIndex = make([]int, len(rf.peers))
					rf.matchIndex = make([]int, len(rf.peers))
					for i := range rf.peers {
						rf.nextIndex[i] = 1
						rf.matchIndex[i] = 0
					}

					go func() {
						i := 0
						for {
							i++
							rf.log = append(rf.log, LogEntry{rf.currentTerm, i, fmt.Sprintf("user send : %d", i)})
							time.Sleep(3 * time.Second)
						}
					}()
				}
			case Leader:
				rf.broadcastHeartbeat()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

type VoteArgs struct {
	Term        int
	CandidateID int
}

type VoteReply struct {
	Term        int
	VoteGranted bool
}

func (rf *Raft) broadcastRequestVote() {
	var args = VoteArgs{
		Term:        rf.currentTerm,
		CandidateID: rf.id,
	}

	for i := range rf.peers {
		go func(i int) {
			var reply VoteReply
			rf.sendRequestVote(i, args, &reply)
		}(i)
	}
}

func (rf *Raft) sendRequestVote(serverID int, args VoteArgs, reply *VoteReply) {
	client, err := rpc.DialHTTP("tcp", rf.peers[serverID].address)
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	defer client.Close()
	client.Call("Raft.RequestVote", args, reply)

	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.state = Follower
		rf.votedFor = -1
		return
	}

	if reply.VoteGranted {
		rf.voteCount++
	}

	if rf.voteCount >= len(rf.peers)/2+1 {
		rf.leaderC <- true
	}
}

type HeartbeatArgs struct {
	Term     int
	LeaderID int

	PrevLogIndex int
	PrevLogTerm  int

	Entries      []LogEntry
	LeaderCommit int
}

type HeartbeatReply struct {
	Success   bool
	Term      int
	NextIndex int
}

func (rf *Raft) broadcastHeartbeat() {
	for i := range rf.peers {

		var args HeartbeatArgs
		args.Term = rf.currentTerm
		args.LeaderID = rf.id
		args.LeaderCommit = rf.commitIndex

		prevLogIndex := rf.nextIndex[i] - 1
		if rf.getLastIndex() > prevLogIndex {
			args.PrevLogIndex = prevLogIndex
			args.PrevLogTerm = rf.log[prevLogIndex].LogTerm
			args.Entries = rf.log[prevLogIndex:]
			log.Printf("send entries: %v\n", args.Entries)
		}

		go func(i int, args HeartbeatArgs) {
			var reply HeartbeatReply
			rf.sendHeartbeat(i, args, &reply)
		}(i, args)
	}
}

func (rf *Raft) sendHeartbeat(serverID int, args HeartbeatArgs, reply *HeartbeatReply) {
	client, err := rpc.DialHTTP("tcp", rf.peers[serverID].address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	defer client.Close()
	client.Call("Raft.Heartbeat", args, reply)

	if reply.Success {
		if reply.NextIndex > 0 {
			rf.nextIndex[serverID] = reply.NextIndex
			rf.matchIndex[serverID] = rf.nextIndex[serverID] - 1
		}
	} else {
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = Follower
			rf.votedFor = -1
			return
		}
	}
}

func (rf *Raft) getLastIndex() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogIndex
}

func (rf *Raft) getLastTerm() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogTerm
}
