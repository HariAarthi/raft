package shared

type RaftConsensus interface {

	IncomingVoteRequest(args *RequestVoteArgs, reply *RequestVoteReply) error
	
	AppendEntries(args *RequestVoteArgs, reply *RequestVoteReply) error

}