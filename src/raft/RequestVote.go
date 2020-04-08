package raft

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).-----------------------------------
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).-----------------------------------------
	Term        int
	VoteGranted bool
}

//
// example RequestVote RPC handler.handler、handler、handler
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).----------------------------------------
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		if args.LastLogIndex >= rf.commitIndex && rf.log[rf.commitIndex].Term==args.LastLogTerm{
			rf.voteFor = args.CandidateId
			rf.findBiggerChan <- 1
			DPrintf("%d votefor %d,当前 term %d", rf.me, args.CandidateId, rf.currentTerm)
			reply.Term = args.Term
			reply.VoteGranted = true
			return 
		}
		if args.LastLogIndex >= rf.commitIndex && rf.log[rf.commitIndex].Term!=args.LastLogTerm{
			DPrintf("%d what ? what the fuck?!", rf.me)
			reply.Term = args.Term
			reply.VoteGranted = false
			return 
		}
		DPrintf("%d leader %d的lastcommit害没自己的 commit大", rf.me,args.CandidateId)
		reply.Term = args.Term
		reply.VoteGranted = false
	} else {
		reply.Term = args.Term
		reply.VoteGranted = false
	}
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}
