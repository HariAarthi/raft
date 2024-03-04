package server

import (
 //"net/rpc"
 //"time"
 "fmt"
)
type ServerNode struct {

	serverId int
	
	consensusModule *RaftNode
	
	exitChannel chan struct{}

}

func NewServer(serverId int, exitChannel chan struct{}) *ServerNode {

	server := new(ServerNode)
	server.serverId = serverId
	server.consensusModule = NewConsensusModule(serverId)
	server.exitChannel = exitChannel
	
	return server
}

func (server *ServerNode) SpinUp() {
	server.consensusModule.SpinUp();
	fmt.Println("Server spinup complete")	
}