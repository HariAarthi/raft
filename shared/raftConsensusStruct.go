package shared

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntries struct {
	Term        int
	leaderId    int
	prevLogIndex int
	prevLogTerm int
	leaderCommit int
}

type AppendEntriesReply struct {
	Term        int
	Success     bool
}