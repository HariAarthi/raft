package server

import (
	"fmt"
	"time"
	"net"
	"net/rpc"
	"log"
	"raft/shared"
	"sync"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
	Dead
)

func (state State) String() string {
	switch state {
		case Follower:
			return "Follower"
		case Candidate:
			return "Candidate"
		case Leader:
			return "Leader"
		case Dead:
			return "Dead"
		default:
			panic("unreachable")
	}

}

type RaftConsensus int

type RaftNode struct {

	mutex sync.Mutex

	requestVoteChannel chan struct{}
	serverId int
	serverPort string
	
	rpcServer *rpc.Server
	peerClientPorts []string
	peerClients []*rpc.Client
	
	raftConsensus RaftConsensus
	
	nodeState State
	
	currentTerm int
	votedFor int
	
	electionResetEvent time.Time
}

func NewConsensusModule(serverId int) (*RaftNode) {
	
	ports := getPorts()
	
	raftNode := new(RaftNode)
	raftNode.requestVoteChannel = make(chan struct{})
	raftNode.serverId = serverId
	raftNode.serverPort = ports[serverId]
	
	log.Print("listening on port ", raftNode.serverPort)
	
	raftNode.SpinUp()
	
	//raftNode.raftConsensus := new(RaftConsensus)
	
	raftNode.rpcServer = rpc.NewServer()
	raftNode.rpcServer.RegisterName("RaftNode",raftNode)
	listener, _ := net.Listen("tcp", raftNode.serverPort)
	raftNode.rpcServer.Accept(listener)
	
	
	raftNode.initializeRpcClients(ports)
	log.Print ("initialized ConsensusModule...")
	log.Print ("listener... ",listener)
	return raftNode
}

func (raftNode *RaftNode) SpinUp() {
	go func() {
		raftNode.electionResetEvent = time.Now()
		raftNode.nodeState = Candidate
		raftNode.initializeRpcClients(getPorts())
		raftNode.initiateElectionTimer()
	} ()
}

func (raftNode *RaftNode) initiateElectionTimer() {
	go func () {
		ticker := time.NewTicker(1000 * time.Millisecond)
		timeOutDuration := time.Duration(5 * ((raftNode.serverId) + 1) ) * time.Second
		for {
			//time.Sleep(time.Duration(5 * ((raftNode.serverId) + 1) ) * time.Second)
			
			<-ticker.C
			
			if elapsed := time.Since(raftNode.electionResetEvent); elapsed >= timeOutDuration {
				fmt.Println("state of raftNode...", raftNode.nodeState)
				//if(raftNode.nodeState == Candidate) {
					fmt.Println("Now, its time to ask for vote...", raftNode.serverId)	
					raftNode.conductElection()
				//}
			}
		}
	} ()
}

func (raftNode *RaftNode) conductElection() {

	raftNode.nodeState = Candidate
	raftNode.currentTerm += 1
	raftNode.electionResetEvent = time.Now()
	raftNode.votedFor = raftNode.serverId
	
	savedCurrentTerm := raftNode.currentTerm
	
	fmt.Println("count of clients ", len(raftNode.peerClients))
	if len(raftNode.peerClients) == 0 {	
		raftNode.initializeRpcClients(getPorts())
	}
	
	if len(raftNode.peerClients) > 0 {	
				votesReceived := 1
				for index, rpcClient := range raftNode.peerClients {
				
					args := shared.RequestVoteArgs{
						Term:        1,
						CandidateId: 1,
					}
					var reply shared.RequestVoteReply
					
					if rpcClient != nil {
						fmt.Println("Contacting the port for vote ", index)
						log.Print("Calling IncomingVoteRequest ")
						err := rpcClient.Call("RaftNode.IncomingVoteRequest", args, &reply)
						if err != nil {
							fmt.Println("error while calling IncomingVoteRequest ", err.Error())
						} else if reply.Term > savedCurrentTerm {
							fmt.Println("term out of date in RequestVoteReply")
							//raftNode.becomeFollower(reply.Term)
							return
						} else if reply.Term == savedCurrentTerm {
							if reply.VoteGranted {
								votesReceived += 1
								if votesReceived*2 > len(raftNode.peerClients)+1 {
									// Won the election!
									fmt.Println("wins election with %d votes", votesReceived)
									raftNode.startLeader()
									return
								}
							}
						}
					}
				}
	}
	
}

func (raftNode *RaftNode) startLeader() error {

	log.Print("start leader")
		
	return nil
}

func (raftNode *RaftNode) IncomingVoteRequest(args *shared.RequestVoteArgs, reply *shared.RequestVoteReply) error {

	log.Print("got an incoming vote request")
	
	if args.Term > raftNode.currentTerm {
		raftNode.nodeState = Follower
		raftNode.currentTerm = args.Term
	} else if args.Term == raftNode.currentTerm {
		raftNode.nodeState = Candidate
	}
	
	raftNode.electionResetEvent = time.Now()
	
	return nil
}

func (raftNode *RaftNode) AppendEntries(args *shared.AppendEntries, reply *shared.AppendEntriesReply) error {

	log.Print("got an append entries request")
	
	if(args.Term > raftNode.currentTerm) {
		raftNode.nodeState = Follower
	}
	
	if(raftNode.currentTerm > args.Term) {
		reply.Success = false
	}
	
	raftNode.electionResetEvent = time.Now()
	reply.Term = raftNode.currentTerm
	
	return nil
}

func (raftNode *RaftNode) initializeRpcClients(ports []string) {
	
	ports = getPorts()
	log.Print("method initializeRpcClients ")
	for index, portNumber := range ports {
		if index != raftNode.serverId {
			raftNode.peerClientPorts = append(raftNode.peerClientPorts, portNumber)
			rpcClient, err := rpc.Dial("tcp",portNumber)
			if err != nil {
				//log.Print("error during connection: ", err)
			}
			if err == nil {
				log.Print("rpcClient for port ", rpcClient, portNumber)
				raftNode.peerClients = append(raftNode.peerClients, rpcClient)
			}
		}
		index = index + 1
	}
}

/*
func (raftConsensus *RaftConsensus) IncomingVoteRequest(args *shared.RequestVoteArgs, reply *shared.RequestVoteReply) error {
	fmt.Println("incoming vote request")
	
	return nil
}
*/

func getPorts() ([]string) {
	var ports = []string {"localhost:2222","localhost:2223","localhost:2224","localhost:2225"}
	return ports
}