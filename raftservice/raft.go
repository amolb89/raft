package raftservice

import (
	"encoding/json"
	clust "github.com/amolb89/raft/clust"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	REQUESTVOTE = 1
)
const (
	APEENDENTRIES = 2
)
const (
	OK = 3
)
const (
	RESPONSE_APPENDENTRIES = 4
)

type Raft interface {
	Term() int
	isLeader() bool
	Outbox() <-chan *clust.LogItem
	Inbox() chan<- string
}

type RaftImpl struct {
	pid         int
	term        int
	isLeader    bool
	timeout     int
	noOfVotes   int
	logFile     string
	commitIndex int64
	lastApplied int64 //index of the highest log entry applied to the state machine
	log         []*clust.LogItem
	LastIndex   int
	outbox      chan<- *clust.LogItem
	inbox       <-chan string
	server      *clust.MockServer
	currLeader  chan int
	timer       *time.Timer
	//For logging purpose
	TRACE *log.Logger
	ERROR *log.Logger
}

type Log struct {
	LogEntries []*clust.LogItem
}

func (r *RaftImpl) Inbox() <-chan string {
	return r.inbox
}

func (r *RaftImpl) Outbox() chan<- *clust.LogItem {
	return r.outbox
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

func (r *RaftImpl) LoadLog(c *Log) error {
	fileBytes, _ := ioutil.ReadFile(r.logFile)
	err := json.Unmarshal(fileBytes, c)
	return err
}

func (r *RaftImpl) AddToLog(logItem *clust.LogItem, index int) {
	if len(r.log) == index {
		r.AppendLogItem(logItem)
	} else {
		r.log = r.log[:index]
		r.AppendLogItem(logItem)
	}
}

func (r *RaftImpl) AppendLogItem(logItem *clust.LogItem) {
	//Add to memory log
	r.log = append(r.log, logItem)
	logBytes, _ := json.Marshal(r.log)
	ioutil.WriteFile(r.logFile, logBytes, 0644)
}

func (r *RaftImpl) SendHeartBeats() {
	if r.isLeader == true {
		r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " is sending AppendEntries heartbeats")
		msg := new(clust.Envelope)
		msg.Pid = -1
		msg.TermId = r.Term()
		msg.Msg = nil
		msg.Resp = nil
		msg.MsgType = APEENDENTRIES
		r.server.Outbox() <- msg
	}
}

func (r *RaftImpl) Configure(servId int, configPath string, currTimeout int, currTerm int, leaderchan chan int, prob float64, delay int, exit <-chan int) {
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	r.TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	r.pid = servId
	r.term = currTerm
	r.isLeader = false
	r.timeout = currTimeout
	r.server = new(clust.MockServer)
	r.server.Set(servId, configPath, prob, delay, exit)
	r.TRACE.Println("Configured Server ", r.Pid(), " with timeout ", r.timeout)
	r.currLeader = leaderchan
	// Read log file and populate log
	r.logFile = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@" + strconv.Itoa(r.Pid())
	log := new(Log)
	r.LoadLog(log)
	r.log = log.LogEntries
	for i := range r.log {
		r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Data :", r.log[i].Data)
	}
	r.LastIndex = len(r.log) - 1
	r.timer = time.NewTimer(time.Duration(r.timeout) * time.Second)
	//tempTimer := time.NewTimer(time.Second * 2)
	//r.timer = tempTimer
	r.TRACE.Println("Loaded log on server ", r.Pid())
	//Create inbox and outbox channels
	r.outbox = make(chan *clust.LogItem, 100)
	r.inbox = make(chan string, 100)
	go r.follower()
}

func (r *RaftImpl) follower() {
	r.TRACE.Println("Server ", r.Pid(), " is in follower state with timeout :", r.timeout)
	r.isLeader = false
	for {
		select {
		case msg := <-r.server.Inbox():
			r.TRACE.Println("Message received at server ", r.Pid(), " with currect term ", r.Term(), " with msgType ", msg.MsgType, " with msgTerm ", msg.TermId)
			if msg.MsgType == REQUESTVOTE && msg.TermId > r.term {
				//Send OK
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received RequestVote from server ", msg.Pid, " with term ", msg.TermId)
				r.TRACE.Println("Server ", r.Pid(), " sending Ok message to server ", msg.Pid)
				r.term = msg.TermId
				okMsg := new(clust.Envelope)
				okMsg.Pid = msg.Pid
				okMsg.TermId = r.Term()
				okMsg.MsgType = OK
				r.server.Outbox() <- okMsg
				//1 ok in one term
			} else if msg.MsgType == REQUESTVOTE && msg.TermId <= r.term {
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " rejected RequestVote from server ", msg.Pid, "with term ", msg.TermId)
				//Reject the request -> stale term
			} else if msg.MsgType == APEENDENTRIES && r.term <= msg.TermId {
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received AppendEntries from ", msg.Pid, " with term ", msg.TermId, " msg is ", msg.Msg, " Complete msg is ", msg)
				if msg.Msg != nil {
					appendEntry := msg.Msg
					responseMsg := new(clust.Envelope)
					responseMsg.Pid = msg.Pid
					responseMsg.TermId = r.Term()
					r.TRACE.Println("AppendEntry received is : appendEntry.prevLogIndex - ", appendEntry.PrevLogIndex, " appendEntry.prevLogTerm - ", appendEntry.PrevLogTerm, " LeaderCommit - ", appendEntry.LeaderCommit)
					var response *clust.Response
					if appendEntry.PrevLogIndex == -1 {
						response = new(clust.Response)
						//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
						response.MatchIndex = 0
						//response.MatchTerm = appendEntry.PrevLogTerm
						response.Success = true
						//Append Entry
						if appendEntry.LogItm != nil {
							r.TRACE.Println("Server ", r.Pid(), " added log entry with index 0 and term ", appendEntry.LogItm.Term, " and data ", appendEntry.LogItm.Data)
							r.AddToLog(appendEntry.LogItm, 0)
							r.LastIndex = 0
						}
					} else if appendEntry.PrevLogIndex > r.LastIndex {
						response = new(clust.Response)
						//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
						response.MatchIndex = appendEntry.PrevLogIndex
						//response.MatchTerm = appendEntry.PrevLogTerm
						response.Success = false
						r.TRACE.Println("Server pid:", r.Pid(), " No match for log entries at index ", appendEntry.PrevLogIndex)
					} else if r.log[appendEntry.PrevLogIndex].Term == appendEntry.PrevLogTerm {
						response = new(clust.Response)
						//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
						response.MatchIndex = appendEntry.PrevLogIndex
						//response.MatchTerm = appendEntry.PrevLogTerm
						response.Success = true
						r.TRACE.Println("Terms match on index ", appendEntry.PrevLogIndex)
						//Append Entry
						if appendEntry.LogItm != nil {
							response.MatchIndex = appendEntry.PrevLogIndex + 1
							//response.MatchTerm = appendEntry.PrevLogTerm
							r.TRACE.Println("Server ", r.Pid(), " added log entry with index ", (appendEntry.PrevLogIndex + 1), " and term ", appendEntry.LogItm.Term, " and data ", appendEntry.LogItm.Data)
							r.AddToLog(appendEntry.LogItm, appendEntry.PrevLogIndex+1)
							r.LastIndex = r.LastIndex + 1
						}
					} else {
						response = new(clust.Response)
						//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
						response.MatchIndex = appendEntry.PrevLogIndex
						//response.MatchTerm = appendEntry.PrevLogTerm
						response.Success = false
						r.TRACE.Println("Server pid:", r.Pid(), " No match for log entries at index ", appendEntry.PrevLogIndex)
					}
					responseMsg.MsgType = RESPONSE_APPENDENTRIES
					responseMsg.Msg = nil
					responseMsg.Resp = response
					r.server.Outbox() <- responseMsg
				} else {
					r.TRACE.Println("Received heartbeat message ")
					r.TRACE.Println("Log for server ", r.Pid())
					for _, value := range r.log {
						r.TRACE.Print(" ", value)
					}
				}
				r.term = msg.TermId
				r.timer.Reset(time.Duration(r.timeout) * time.Second)
			}
		case <-r.timer.C:
			r.TRACE.Println("Timeout in follower state:: Server ", r.Pid())
			go r.candidate()
			return
		}
	}
}

func (r *RaftImpl) candidate() {
	r.TRACE.Println("Server ", r.Pid(), " is in candidate state")
	r.term = r.Term() + 1
	// Send RequestVote broadcast
	r.TRACE.Println("Server ", r.Pid(), " broadcasting RequestVote message")
	msg := new(clust.Envelope)
	msg.Pid = -1
	msg.TermId = r.Term()
	msg.MsgType = REQUESTVOTE
	r.TRACE.Println("Server ", r.Pid(), " broadcasting RequestVote message with msgType ", msg.MsgType)
	r.server.Outbox() <- msg
	//Reset election timeout
	r.timeout = clust.Random(5, 30)
	r.timer.Reset(time.Duration(r.timeout) * time.Second)
	r.TRACE.Println("Server ::Candidate State::", r.Pid(), " timeout reset to : ", r.timeout)
	r.noOfVotes = 1

	for {
		select {
		case msg := <-r.server.Inbox():
			switch {
			case msg.MsgType == OK:
				r.noOfVotes = r.noOfVotes + 1
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received Ok message from server ", msg.Pid, " with term ", msg.TermId, " current number of votes ", r.noOfVotes, " and number of peers is ", r.server.NoOfPeers(), " Required number of votes are ", (r.server.NoOfPeers()+1)/2)
				if r.noOfVotes > (r.server.NoOfPeers()+1)/2 {
					r.noOfVotes = 0
					go r.leader()
					return
				}
			case msg.MsgType == REQUESTVOTE:
				if msg.TermId > r.term {
					r.term = msg.TermId
					okMsg := new(clust.Envelope)
					okMsg.Pid = msg.Pid
					okMsg.TermId = r.Term()
					okMsg.MsgType = OK
					r.server.Outbox() <- okMsg
					go r.follower()
					return
				}
			case msg.MsgType == APEENDENTRIES:
				if msg.TermId > r.term {
					r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received AppendEntries from ", msg.Pid, " with term ", msg.TermId)
					if msg.Msg != nil {
						appendEntry := msg.Msg
						responseMsg := new(clust.Envelope)
						responseMsg.Pid = msg.Pid
						responseMsg.TermId = r.Term()
						r.TRACE.Println("AppendEntry received is : appendEntry.prevLogIndex - ", appendEntry.PrevLogIndex, " appendEntry.prevLogTerm - ", appendEntry.PrevLogTerm, " LeaderCommit - ", appendEntry.LeaderCommit)
						var response *clust.Response
						if appendEntry.PrevLogIndex == -1 {
							response = new(clust.Response)
							//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
							response.MatchIndex = 0
							//response.MatchTerm = appendEntry.PrevLogTerm
							response.Success = true
							//Append Entry
							if appendEntry.LogItm != nil {
								r.TRACE.Println("Server ", r.Pid(), " added log entry with index 0 and term ", appendEntry.LogItm.Term, " and data ", appendEntry.LogItm.Data)
								r.AddToLog(appendEntry.LogItm, 0)
								r.LastIndex = 0
							}
						} else if r.log[appendEntry.PrevLogIndex].Term == appendEntry.PrevLogTerm {
							response = new(clust.Response)
							//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
							response.MatchIndex = appendEntry.PrevLogIndex
							//response.MatchTerm = appendEntry.PrevLogTerm
							response.Success = true
							r.TRACE.Println("Terms match on index ", appendEntry.PrevLogIndex)
							//Append Entry
							if appendEntry.LogItm != nil {
								response.MatchIndex = appendEntry.PrevLogIndex + 1
								//response.MatchTerm = appendEntry.PrevLogTerm
								r.TRACE.Println("Server ", r.Pid(), " added log entry with index ", (appendEntry.PrevLogIndex + 1), " and term ", appendEntry.LogItm.Term, " and data ", appendEntry.LogItm.Data)
								r.AddToLog(appendEntry.LogItm, appendEntry.PrevLogIndex+1)
								r.LastIndex = r.LastIndex + 1
							}
						} else {
							response = new(clust.Response)
							//response = Response{appendEntry.PrevLogIndex,appendEntry.PrevLogTerm,true}
							response.MatchIndex = appendEntry.PrevLogIndex
							//response.MatchTerm = appendEntry.PrevLogTerm
							response.Success = false
							r.TRACE.Println("Server pid:", r.Pid(), " No match for log entries at index ", appendEntry.PrevLogIndex)
						}
						responseMsg.MsgType = RESPONSE_APPENDENTRIES
						responseMsg.Msg = nil
						responseMsg.Resp = response
						r.server.Outbox() <- responseMsg

					}
					r.term = msg.TermId
					r.timer.Reset(time.Duration(r.timeout) * time.Second)
					go r.follower()
					return
				}
				//else it is a stale message which popped up.
			}
		case <-r.timer.C:
			r.TRACE.Println("Timeout in candidate state:: Server ", r.Pid())
			go r.candidate()
			return
		}
	}
}

func (r *RaftImpl) leader() {
	r.TRACE.Println("Server ", r.Pid(), " is in leader state")
	r.isLeader = true
	//Send heartbeat
	r.SendHeartBeats()
	//Volatile states on leader
	var nextIndex map[int]int  //For each server, index of the next log entry to send to that server.
	var matchIndex map[int]int //For each server, index of highest log entry known to be replicated on server.

	nextIndex = make(map[int]int)
	matchIndex = make(map[int]int)

	//Initialize nextIndex and matchIndex
	for id, _ := range r.server.Peers() {
		nextIndex[id] = len(r.log)
		matchIndex[id] = -1
		r.TRACE.Println("nextId[", id, "] = ", nextIndex[id], " matchIndex[", id, "] = ", matchIndex[id])
	}

	r.timer.Reset(time.Duration(5) * time.Second)
	//r.isLeader = true
	r.CurrLeader() <- r.Pid()
	for {
		select {
		case msg := <-r.server.Inbox():
			switch {
			case msg.MsgType == REQUESTVOTE:
				if msg.TermId > r.term {
					r.term = msg.TermId
					//send Ok message
					okMsg := new(clust.Envelope)
					okMsg.Pid = msg.Pid
					okMsg.TermId = r.Term()
					okMsg.MsgType = OK
					r.server.Outbox() <- okMsg
					go r.follower()
					return
				}
				//else it is a stale message
			case msg.MsgType == APEENDENTRIES:
				if msg.TermId > r.term {
					go r.follower()
					return
				}
				//else it is a stale message
			case msg.MsgType == RESPONSE_APPENDENTRIES:
				response := msg.Resp
				if response.Success == true {
					r.TRACE.Println("Success! matchIndex is ", response.MatchIndex)
					if response.MatchIndex <= r.LastIndex {
						//Check to ensure delayed response doesn't cause a problem
						matchIndex[msg.Pid] = response.MatchIndex
						if response.MatchIndex == nextIndex[msg.Pid] {
							nextIndex[msg.Pid] = nextIndex[msg.Pid] + 1
						}
					}

				} else {
					if response.MatchIndex == nextIndex[msg.Pid]-1 {
						nextIndex[msg.Pid] = nextIndex[msg.Pid] - 1
					}
				}
				r.TRACE.Println("New next and matching entries")
				for id, _ := range r.server.Peers() {
					r.TRACE.Println("Server id: ", id, " NextIndex : ", nextIndex[id], " matchIndex : ", matchIndex[id])
				}
			}
			//Message on Inbox to send commands of any kind, and to have them replicated by raft on majority of the server.
		case cmd := <-r.Inbox():
			logItem := new(clust.LogItem)
			logItem.Term = r.Term()
			logItem.Data = cmd
			r.AddToLog(logItem, r.LastIndex)
			r.LastIndex = r.LastIndex + 1
			r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " is sending AppendEntries message")
			for id, _ := range r.server.Peers() {
				if id != r.Pid() {
					msg := new(clust.Envelope)
					msg.Pid = id
					msg.TermId = r.Term()
					msg.MsgType = APEENDENTRIES
					appendEntry := new(clust.AppendEntry)
					appendEntry.PrevLogIndex = nextIndex[id] - 1
					if appendEntry.PrevLogIndex >= 0 {
						appendEntry.PrevLogTerm = r.log[nextIndex[id]-1].Term
					} else {
						appendEntry.PrevLogTerm = 0
					}
					if nextIndex[id] == matchIndex[id]+1 {
						r.TRACE.Println("Sending heartbeat AppendEntry to ", id)
						appendEntry = nil
					} else if nextIndex[id] == r.LastIndex+1 {
						r.TRACE.Println("Sending AppendEntry to ", id, " with nil log item at index ", nextIndex[id])
						appendEntry.LogItm = nil
						appendEntry.LeaderCommit = 1
					} else {
						r.TRACE.Println("Sending AppendEntry to ", id, " with log item at index ", nextIndex[id])
						appendEntry.LogItm = r.log[nextIndex[id]]
						appendEntry.LeaderCommit = 1
					}
					msg.Msg = appendEntry
					msg.Resp = nil
					r.server.Outbox() <- msg
				}
			}
		case <-r.timer.C:
			r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " is sending AppendEntries message")
			for id, _ := range r.server.Peers() {
				if id != r.Pid() {
					msg := new(clust.Envelope)
					msg.Pid = id
					msg.TermId = r.Term()
					msg.MsgType = APEENDENTRIES
					appendEntry := new(clust.AppendEntry)
					appendEntry.PrevLogIndex = nextIndex[id] - 1
					if appendEntry.PrevLogIndex >= 0 {
						appendEntry.PrevLogTerm = r.log[nextIndex[id]-1].Term
					} else {
						appendEntry.PrevLogTerm = 0
					}
					if r.LastIndex == matchIndex[id] {
						r.TRACE.Println("Sending heartbeat AppendEntry to ", id)
						appendEntry = nil
					} else if nextIndex[id] == r.LastIndex+1 {
						r.TRACE.Println("Sending AppendEntry to ", id, " with nil log item at index ", nextIndex[id])
						appendEntry.LogItm = nil
						appendEntry.LeaderCommit = 1
					} else {
						r.TRACE.Println("Sending AppendEntry to ", id, " with log item at index ", nextIndex[id])
						appendEntry.LogItm = r.log[nextIndex[id]]
						appendEntry.LeaderCommit = 1
					}
					msg.Msg = appendEntry
					msg.Resp = nil
					r.server.Outbox() <- msg
				}
			}
			r.timer.Reset(time.Duration(5) * time.Second)
		}
	}
}
