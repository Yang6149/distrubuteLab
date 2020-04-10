package raft

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PreLogIndex  int
	PreLogTerm   int
	Entries      []Entry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term       int
	Success    bool
	MatchIndex int
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term >= rf.currentTerm {
		rf.state = follower
		rf.currentTerm = args.Term
		rf.appendChan <- 1
		DPrintf("%d：重置ele时间", rf.me)
		if args.PreLogIndex >= len(rf.log) {
			//preIndex 越界
			DPrintf("%d :因为leader传来的preIndex大于自己的len(log)所以return false", rf.me)
			reply.Term = rf.currentTerm
			reply.Success = false
			return
		}
		if rf.log[args.PreLogIndex].Term != args.PreLogTerm {

			//index and  term can't match,return false
			reply.Term = rf.currentTerm
			reply.Success = false
			DPrintf("%d :index and term can't match ,return false,my last Command is ", rf.me, rf.log[args.PreLogIndex].Command)
		} else {
			// index and term is matched
			reply.Term = rf.currentTerm
			reply.Success = true
			if len(args.Entries) > 0 {

				//sent a entry
				//determine whether it is overwrited
				if args.PreLogIndex+1 < len(rf.log) {
					//如果我要 overwrite 的内容和我要写的内容一模一样，就证明有一个线程走在我前面，我退出
					if args.Entries[0].Term == rf.log[args.PreLogIndex+1].Term {
						reply.MatchIndex = -1
						DPrintf("%d我这个是个重复overwrite操作，有个线程在我之前，直接退出", rf.me)
						return
					}
					//if it is overwrite
					DPrintf("%d 原来的 log 为 %d ", rf.me, rf.log)
					rf.log[args.PreLogIndex+1] = args.Entries[0]
					rf.log = rf.log[0 : args.PreLogIndex+2]
					DPrintf("%d overwrite 后的 log 为%d", rf.me, rf.log)
				} else {
					//else append to log
					rf.log = append(rf.log, args.Entries[0])
					DPrintf("%d 添加一个新log,现在长度为：%d", rf.me, len(rf.log))
				}
				DPrintf("%d ：现在的log长度是 %d", rf.me, len(rf.log))
				reply.MatchIndex = args.PreLogIndex + 1
				DPrintf("%d ：成功接收一个 entry，matchIndex 为%d ,the entry is ", rf.me, reply.MatchIndex, args.Entries)
			} else {
				//just a heartbeat
				reply.Term = rf.currentTerm
				reply.Success = true
				reply.MatchIndex = args.PreLogIndex
				DPrintf("%d ：just a heartBeat 为%d", rf.me, reply.MatchIndex)
			}
			if args.LeaderCommit > rf.commitIndex {
				//commit all index before args.LeaderCommit
				newCommitNum := min(args.LeaderCommit, reply.MatchIndex)
				if newCommitNum > rf.commitIndex {
					DPrintf("%dcommit index %d", rf.me, newCommitNum)
					rf.commitIndex = newCommitNum
					rf.sendApply <- rf.commitIndex
					DPrintf("%d commit index %d finished", rf.me, rf.commitIndex)
				}
			}
		}

	} else {
		//我的 Term 更大 返回false
		DPrintf("%d :因为自己的term比leader %d 更大return false", rf.me, args.LeaderId)
		reply.Term = rf.currentTerm
		reply.Success = false
	}
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

func min(x int, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}
