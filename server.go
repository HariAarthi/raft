package main

import(
	"fmt"
	//"net/rpc"
	//"log"
	"raft/server"
	"os"
	"strconv"
)

func main() {	
	
	args := os.Args
	var serverId int
	
	if len(args) > 1 {
		serverId, _ = strconv.Atoi(args[1])
		
	} else {
		panic("Server Id missing")
	}
	
	fmt.Println("serverId ", serverId)
	exit := make(chan struct{})
	server.NewServer(serverId, exit).SpinUp()
	
	fmt.Println("waiting to exit...")
	
	<-exit
	
}