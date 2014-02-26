package raftservice

import (
	clust "github.com/amolb89/raft/clust"
	"log"
	"os"
	"time"
)

type Raft interface {
	Term() int
	isLeader() bool
}

type RaftImpl struct {
	pid        int
	term       int
	isLeader   bool
	timeout    int
	noOfVotes  int
	server     *clust.MockServer
	currLeader chan int
	//For logging purpose
	TRACE *log.Logger
	ERROR *log.Logger
}

func (r *RaftImpl) Pid() int {
	return r.pid
}

func (r *RaftImpl) Term() int {
	return r.term
}

func (r *RaftImpl) IsLeader() bool {
	return r.isLeader
}

func (r *RaftImpl) CurrLeader() chan int {
	return r.currLeader
}

func (r *RaftImpl) Configure(servId int, configPath string, currTimeout int, currTerm int, leaderchan chan int, prob float64, delay int) {
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	r.TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	r.pid = servId
	r.term = currTerm
	r.isLeader = false
	r.timeout = currTimeout
	r.server = new(clust.MockServer)
	r.server.Set(servId, configPath, prob, delay)
	r.TRACE.Println("Configured Server ", r.Pid(), " with timeout ", r.timeout)
	r.currLeader = leaderchan
	go r.follower()
}

func (r *RaftImpl) follower() {
	r.TRACE.Println("Server ", r.Pid(), " is in follower state")
	r.isLeader = false
	nextTerm := true
L:
	for {
		select {
		case msg := <-r.server.Inbox():
			if msg.Msg == "RequestVote" && nextTerm == true {
				//Send OK
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received RequestVote from server ", msg.Pid, " with term ", msg.TermId)
				if r.Term() <= msg.TermId {
					r.TRACE.Println("Server ", r.Pid(), " sending Ok message to server ", msg.Pid)
					r.term = msg.TermId
					okMsg := new(clust.Envelope)
					okMsg.Pid = msg.Pid
					okMsg.TermId = r.Term()
					okMsg.Msg = "Ok"
					r.server.Outbox() <- okMsg
					//1 ok in one term
					nextTerm = false
				} else {
					//Transition to candidate state
					break L
				}
			} else if msg.Msg == "AppendEntries" {
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received AppendEntries from ", msg.Pid, " with term ", msg.TermId)
				r.term = msg.TermId
				nextTerm = true
			}
			//Else it is AppendEntries -> No action needed
		case <-time.After(time.Duration(r.timeout) * time.Second):
			r.TRACE.Println("Timeout in follower state:: Server ", r.Pid())
			break L
		}
	}
	go r.candidate()
}

func (r *RaftImpl) candidate() {
	r.TRACE.Println("Server ", r.Pid(), " is in candidate state")
	// Send RequestVote broadcast
	r.TRACE.Println("Server ", r.Pid(), " broadcasting RequestVote message")
	msg := new(clust.Envelope)
	msg.Pid = -1
	msg.TermId = r.Term()
	msg.Msg = "RequestVote"
	r.term = r.Term() + 1
	r.server.Outbox() <- msg
	//Reset election timeout
	r.timeout = clust.Random(5, 30)
	r.TRACE.Println("Server ", r.Pid(), " timeout reset to : ", r.timeout)
	r.noOfVotes = 1
L:
	for {
		select {
		case msg := <-r.server.Inbox():
			switch {
			case msg.Msg == "Ok":
				r.noOfVotes = r.noOfVotes + 1
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received Ok message from server ", msg.Pid, " with term ", msg.TermId, " current number of votes ", r.noOfVotes, " and number of peers is ", r.server.NoOfPeers(), " Required number of votes are ", (r.server.NoOfPeers()+1)/2)
				if r.noOfVotes > (r.server.NoOfPeers()+1)/2 {
					r.noOfVotes = 0
					go r.leader()
					break L
				}
			//ignore any requestvote as 1 vote in one term
			case msg.Msg == "AppendEntries":
				r.TRACE.Println("Server ", r.Pid(), " in candidate state received AppendEntries from server ", msg.Pid)
				r.noOfVotes = 0
				go r.follower()
				break L
			}
		case <-time.After(time.Duration(r.timeout) * time.Second):
			r.TRACE.Println("Timeout in candidate state:: Server ", r.Pid())
			go r.follower()
			break L
		}
	}
}

func (r *RaftImpl) leader() {
	r.TRACE.Println("Server ", r.Pid(), " is in leader state")
	r.isLeader = true
	r.CurrLeader() <- r.Pid()
	for {
		select {
		case <-time.After(5):
			r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " is sending AppendEntries message")
			msg := new(clust.Envelope)
			msg.Pid = -1
			msg.TermId = r.Term()
			msg.Msg = "AppendEntries"
			r.server.Outbox() <- msg
			r.term = r.Term() + 1
		}
	}
}
