package raftservice

import (
	"encoding/json"
	clust "github.com/amolb89/raft/clust"
	levigo "github.com/jmhodges/levigo"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	"sync"
	"sort"
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
	COMMAND = 5
)

const (
	GET = 1
	PUT = 2
	DELETE = 3
)

type Raft interface {
	Term() int
	isLeader() bool
	Outbox() <-chan *clust.LogItem
	Inbox() chan<- string
}

type RaftImpl struct {
	pid         int32
	leaderPid   int32
	term        int32
	isLeader    bool
	timeout     int
	noOfVotes   int
	logFile     string
	logDB	    *levigo.DB
	database    string
	termFile    string
	commitIndex int
	db	    *levigo.DB
	ro 	    *levigo.ReadOptions
	wo 	    *levigo.WriteOptions
	lastApplied int //index of the highest log entry applied to the state machine
	log         []*clust.LogItem
	LastIndex   int
	outbox      chan *clust.LogItem
	inbox       chan *clust.Command
	server      *clust.MockServer
	currLeader  chan int32
	timer       *time.Timer
	//For logging purpose
	TRACE *log.Logger
	ERROR *log.Logger
}

type Log struct {
	LogEntries []*clust.LogItem
}

func (r *RaftImpl) Server() *clust.MockServer {
	return r.server
}

func (r *RaftImpl) Inbox() chan *clust.Command {
	return r.inbox
}

func (r *RaftImpl) Outbox() chan *clust.LogItem {
	return r.outbox
}

func (r *RaftImpl) Pid() int32 {
	return r.pid
}

func (r *RaftImpl) Term() int32 {
	return r.term
}

func (r *RaftImpl) IsLeader() bool {
	return r.isLeader
}

func (r *RaftImpl) CurrLeader() chan int32 {
	return r.currLeader
}

func (r *RaftImpl) LoadLog() {
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	var err error
	r.logDB, err = levigo.Open(r.logFile, opts)
	if err != nil {
		r.TRACE.Println(err)
	}
	
	it := r.logDB.NewIterator(r.ro)
	defer it.Close()
		
	it.SeekToFirst()
	var index int32
	for it=it; it.Valid(); it.Next() {
		l := new(clust.LogItem)
		err := json.Unmarshal(it.Value(), l)
		if err != nil {
			r.TRACE.Println(err)
		} 
		err = json.Unmarshal(it.Key(), &index)
		if err != nil {
			r.TRACE.Println(err)
		} 
		if int(index) <= cap(r.log) && int(index) > len(r.log) {
			r.log = r.log[:index]
			r.log[index-1]  = l
		} else if int(index) > cap(r.log) {
			t := make([]*clust.LogItem, len(r.log), index)
			copy(t, r.log)
			r.log = t
			r.log[index-1] = l
		} else if int(index) <= len(r.log) {
			r.log[index-1] = l
		}
	}
}

func (r *RaftImpl) changeTerm(newTerm int32) {
	termBytes, err := json.Marshal(newTerm)
	if err != nil {
		r.TRACE.Println("Error while changing term ",newTerm)
	}
	err = ioutil.WriteFile(r.termFile, termBytes, 0644)
	if err != nil {
		r.TRACE.Println("Error in updating term")	
	} else {
		 r.term = newTerm
	}
}

func (r *RaftImpl) LoadTerm() {
	fileBytes, _ := ioutil.ReadFile(r.termFile)
	json.Unmarshal(fileBytes, &r.term)	
	r.TRACE.Println("Term for server ",r.Pid()," is: ",r.Term())
}

func (r *RaftImpl) AddToLog(logItem *clust.LogItem, index int) {
	if r.LastIndex + 1 == index {
		r.AppendLogItem(logItem)
	} else {
		r.log = r.log[:index]
		r.RemoveNextEntries(int32(index + 1))
		r.LastIndex = index - 1
		r.AppendLogItem(logItem)
	}
}

func (r *RaftImpl) RemoveNextEntries(index int32) {
	r.TRACE.Println("RemoveNextEntries@",r.Pid()," start Index is ",(index + 1))
	it := r.logDB.NewIterator(r.ro)
	defer it.Close()
	var start int32
	start = index + 1
	indexBytes, _ := json.Marshal(start)
	it.Seek(indexBytes)
	if it.Valid() {
		for it=it; it.Valid(); it.Seek(indexBytes) {
			r.TRACE.Println("removing ",start," from log of server ",r.Pid())
			err := r.logDB.Delete(r.wo,indexBytes)
			if err != nil {
				r.TRACE.Println(err)
			}
			start = start + 1
			indexBytes, _ = json.Marshal(start)
		}
	} else {
		r.TRACE.Println("Seek not valid")
	}
}

func (r *RaftImpl) AppendLogItem(logItem *clust.LogItem) {
	//Add to memory log
	r.log = append(r.log, logItem)
	logBytes, _ := json.Marshal(logItem)
	indexBytes, _ := json.Marshal(len(r.log))
	r.logDB.Put(r.wo,indexBytes,logBytes)
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

func (r *RaftImpl) Configure(servId int32, configPath string, currTimeout int, leaderchan chan int32, prob float64, delay int, exit <-chan int) {
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	r.TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	r.pid = servId
	r.isLeader = false
	r.timeout = currTimeout
	r.server = new(clust.MockServer)
	r.server.Set(servId, configPath, prob, delay, exit)
	r.TRACE.Println("Configured Server ", r.Pid(), " with timeout ", r.timeout)
	r.currLeader = leaderchan
	r.leaderPid = -1
	// Read log file and populate log
	r.logFile = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@" + strconv.Itoa(int(r.Pid()))
	r.database = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/db@" + strconv.Itoa(int(r.Pid()))
	r.termFile = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/term@" + strconv.Itoa(int(r.Pid()))
	r.ro = levigo.NewReadOptions()
	r.wo = levigo.NewWriteOptions()	
	r.log = make([]*clust.LogItem,0,15)
	r.LoadLog()
	for i := range r.log {
		if r.log[i].Data.Type == PUT {
			r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : PUT(", r.log[i].Data.Key,",",r.log[i].Data.Value,")")
		} else if r.log[i].Data.Type == DELETE {
			r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : DELETE(", r.log[i].Data.Key,")")			
		}
	}
	r.LastIndex = len(r.log) - 1
	r.timer = time.NewTimer(time.Duration(r.timeout) * time.Second)
	r.commitIndex = -1
	r.lastApplied = -1
	//tempTimer := time.NewTimer(time.Second * 2)
	//r.timer = tempTimer
	r.TRACE.Println("Loaded log on server ", r.Pid())
	//Create inbox and outbox channels
	r.outbox = make(chan *clust.LogItem, 100)
	r.inbox = make(chan *clust.Command, 100)

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	r.db, _ = levigo.Open(r.database, opts)
	r.LoadTerm()

	go r.proc_outbox()
	go r.follower()
}


func (r *RaftImpl) follower() {
	var mutex = &sync.Mutex{}
	r.TRACE.Println("Server ", r.Pid(), " is in follower state with timeout :", r.timeout)
	r.isLeader = false
	for {
		select {
		case msg := <-r.server.Inbox():
			r.TRACE.Println("Message received at server ", r.Pid(), " with msgType ", msg.MsgType)
			if msg.MsgType == COMMAND {
				if r.leaderPid == -1 {
					r.TRACE.Println("Leader not yet elected")
					r.server.Outbox() <- msg
				} else {
					cmdMsg := new(clust.Envelope)
					cmdMsg.Pid = r.leaderPid
					cmdMsg.MsgType = COMMAND
					cmdMsg.Cmd = msg.Cmd
					r.server.Outbox() <- msg
				}
			} else if msg.MsgType == REQUESTVOTE && msg.TermId > r.term {
				//Send OK
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received RequestVote from server ", msg.Pid, " with term ", msg.TermId)
				r.TRACE.Println("Server ", r.Pid(), " sending Ok message to server ", msg.Pid)
				r.changeTerm(msg.TermId)
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
				r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received AppendEntries from ", msg.Pid, " with term ", msg.TermId, " msg is ", msg.Msg, " Complete msg is ", msg," Leadercommit is ",msg.LeaderCommit)
				r.commitIndex = msg.LeaderCommit
				r.leaderPid = msg.Pid
				if msg.Msg != nil {
					appendEntry := msg.Msg
					responseMsg := new(clust.Envelope)
					responseMsg.Pid = msg.Pid
					responseMsg.TermId = r.Term()
					r.TRACE.Println("AppendEntry received is : appendEntry.prevLogIndex - ", appendEntry.PrevLogIndex, " appendEntry.prevLogTerm - ", appendEntry.PrevLogTerm, " LeaderCommit - ", msg.LeaderCommit)
					var response *clust.Response
					if appendEntry.PrevLogIndex == -1 {
						response = new(clust.Response)
						response.MatchIndex = 0
						response.Success = true
						//Append Entry
						if appendEntry.LogItm != nil {
							r.TRACE.Println("Server ", r.Pid(), " added log entry with index 0 and term ", appendEntry.LogItm.Term, " and Command type : ",appendEntry.LogItm.Data.Type," Key : ", appendEntry.LogItm.Data.Key)
							r.AddToLog(appendEntry.LogItm, 0)
							r.LastIndex = 0
							if msg.LeaderCommit >= r.LastIndex && r.lastApplied == -1{
								// Apply commands from log to the state machine
								r.TRACE.Println("Apply commands from log to the state machine")
								if r.log[0].Data.Type == PUT {
									err := r.db.Put(r.wo,[]byte(r.log[0].Data.Key), []byte(r.log[0].Data.Value))				
									if err == nil {
										val,_ := r.db.Get(r.ro, []byte(r.log[0].Data.Key))
										r.TRACE.Println("Put successful	of the key ",r.log[0].Data.Key," and inserted value is : ",string(val))
									} else {
										r.TRACE.Println("Put unsuccessful	")
									}
								} else if r.log[0].Data.Type == DELETE {
									r.db.Delete(r.wo, []byte(r.log[0].Data.Key))
								}
								r.lastApplied = 0
							}
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
							r.AddToLog(appendEntry.LogItm, appendEntry.PrevLogIndex+1)
							r.LastIndex = r.LastIndex + 1
							r.TRACE.Println("Server ", r.Pid(), " added log entry with index ", (appendEntry.PrevLogIndex + 1), " and term ", appendEntry.LogItm.Term, " and Command type : ",appendEntry.LogItm.Data.Type," Key : ", appendEntry.LogItm.Data.Key," LastIndex is :",r.LastIndex," length is : ",len(r.log)," lastapplied is ",r.lastApplied)
						} else {
							r.log = r.log[:(appendEntry.PrevLogIndex+1)]
							r.RemoveNextEntries(int32(appendEntry.PrevLogIndex+1))
							r.LastIndex = appendEntry.PrevLogIndex
							r.TRACE.Println("length of log is adjusted to ",len(r.log)," lastindex is ",r.LastIndex)
						}
					} else {
						response = new(clust.Response)
						response.MatchIndex = appendEntry.PrevLogIndex
						response.Success = false
						r.TRACE.Println("Server pid:", r.Pid(), " No match for log entries at index ", appendEntry.PrevLogIndex)
					}
					responseMsg.MsgType = RESPONSE_APPENDENTRIES
					responseMsg.Msg = nil
					responseMsg.Resp = response
					r.server.Outbox() <- responseMsg
				} else {
					r.TRACE.Println("Received heartbeat message ")
				}
				if msg.LeaderCommit >= r.LastIndex && r.lastApplied < r.LastIndex {
					r.TRACE.Println("Apply: Length of log: ",len(r.log)," LastIndex: ",r.LastIndex," LastApplied: ",r.lastApplied)
					for i := r.lastApplied + 1 ; i<= r.LastIndex ; i++ {
						r.Outbox() <- r.log[i]
						r.lastApplied = i
					}
				}
				mutex.Lock()
				r.TRACE.Println("Log for server ", r.Pid()," lastApplied ",r.lastApplied," lastindex : ",r.LastIndex," commitindex ",r.commitIndex)
				for i := range r.log {
					if r.log[i].Data.Type == PUT {
						r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : PUT(", r.log[i].Data.Key,",",r.log[i].Data.Value,")")
					} else if r.log[i].Data.Type == DELETE {
						r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : DELETE(", r.log[i].Data.Key,")")			
					}
				}
				mutex.Unlock()
				r.changeTerm(msg.TermId)
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
	var mutex = &sync.Mutex{}
	r.TRACE.Println("Server ", r.Pid(), " is in candidate state")
	//r.term = r.Term() + 1
	r.changeTerm(r.Term() + 1)
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
	r.leaderPid = -1
	for {
		select {
		case msg := <-r.server.Inbox():
			switch {
			case msg.MsgType == COMMAND:
				r.TRACE.Println("Leader not yet elected")
				r.server.Outbox() <- msg
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
					r.changeTerm(msg.TermId)
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
					r.TRACE.Println("Server ", r.Pid(), " with term ", r.Term(), " received AppendEntries from ", msg.Pid, " with term ", msg.TermId, " msg is ", msg.Msg, " Complete msg is ", msg," Leadercommit is ",msg.LeaderCommit)
					r.commitIndex = msg.LeaderCommit
					r.leaderPid = msg.Pid
					if msg.Msg != nil {
						appendEntry := msg.Msg
						responseMsg := new(clust.Envelope)
						responseMsg.Pid = msg.Pid
						responseMsg.TermId = r.Term()
						r.TRACE.Println("AppendEntry received is : appendEntry.prevLogIndex - ", appendEntry.PrevLogIndex, " appendEntry.prevLogTerm - ", appendEntry.PrevLogTerm, " LeaderCommit - ", msg.LeaderCommit)
						var response *clust.Response
						if appendEntry.PrevLogIndex == -1 {
							response = new(clust.Response)
							response.MatchIndex = 0
							response.Success = true
							//Append Entry
							if appendEntry.LogItm != nil {
								r.TRACE.Println("Server ", r.Pid(), " added log entry with index 0 and term ", appendEntry.LogItm.Term, " and Command type : ",appendEntry.LogItm.Data.Type," Key : ", appendEntry.LogItm.Data.Key)
								r.AddToLog(appendEntry.LogItm, 0)
								r.LastIndex = 0
								if msg.LeaderCommit >= r.LastIndex && r.lastApplied == -1{
									// Apply commands from log to the state machine
									r.TRACE.Println("Apply commands from log to the state machine")
									if r.log[0].Data.Type == PUT {
										err := r.db.Put(r.wo,[]byte(r.log[0].Data.Key), []byte(r.log[0].Data.Value))				
										if err == nil {
											val,_ := r.db.Get(r.ro, []byte(r.log[0].Data.Key))
											r.TRACE.Println("Put successful	of the key ",r.log[0].Data.Key," and inserted value is : ",string(val))
										} else {
											r.TRACE.Println("Put unsuccessful	")
										}
									} else if r.log[0].Data.Type == DELETE {
										r.db.Delete(r.wo, []byte(r.log[0].Data.Key))
									}
									r.lastApplied = 0
								}
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
								r.AddToLog(appendEntry.LogItm, appendEntry.PrevLogIndex+1)
								r.LastIndex = r.LastIndex + 1
								r.TRACE.Println("Server ", r.Pid(), " added log entry with index ", (appendEntry.PrevLogIndex + 1), " and term ", appendEntry.LogItm.Term, " and Command type : ",appendEntry.LogItm.Data.Type," Key : ", appendEntry.LogItm.Data.Key," LastIndex is :",r.LastIndex," length is : ",len(r.log)," lastapplied is ",r.lastApplied)
							} else {
								r.log = r.log[:(appendEntry.PrevLogIndex+1)]
								r.RemoveNextEntries(int32(appendEntry.PrevLogIndex+1))
								r.LastIndex = appendEntry.PrevLogIndex
								r.TRACE.Println("length of log is adjusted to ",len(r.log)," lastindex is ",r.LastIndex)
							}
						} else {
							response = new(clust.Response)
							response.MatchIndex = appendEntry.PrevLogIndex
							response.Success = false
							r.TRACE.Println("Server pid:", r.Pid(), " No match for log entries at index ", appendEntry.PrevLogIndex)
						}
						responseMsg.MsgType = RESPONSE_APPENDENTRIES
						responseMsg.Msg = nil
						responseMsg.Resp = response
						r.server.Outbox() <- responseMsg
					} else {
						r.TRACE.Println("Received heartbeat message ")
					}
					if msg.LeaderCommit >= r.LastIndex && r.lastApplied < r.LastIndex {
						r.TRACE.Println("Apply: Length of log: ",len(r.log)," LastIndex: ",r.LastIndex," LastApplied: ",r.lastApplied)
						for i := r.lastApplied + 1 ; i<= r.LastIndex ; i++ {
							r.Outbox() <- r.log[i]
							r.lastApplied = i
						}
					}
					mutex.Lock()
					r.TRACE.Println("Log for server ", r.Pid()," lastApplied ",r.lastApplied," lastindex : ",r.LastIndex," commitindex ",r.commitIndex)
					for i := range r.log {
						if r.log[i].Data.Type == PUT {
							r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : PUT(", r.log[i].Data.Key,",",r.log[i].Data.Value,")")
						} else if r.log[i].Data.Type == DELETE {
							r.TRACE.Println("log@", r.Pid(), " : Index :", i, " Term :", r.log[i].Term, " Command : DELETE(", r.log[i].Data.Key,")")			
						}
					}
					mutex.Unlock()
					r.changeTerm(msg.TermId)
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
	var nextIndex map[int32]int  //For each server, index of the next log entry to send to that server.
	var matchIndex map[int32]int //For each server, index of highest log entry known to be replicated on server.

	nextIndex = make(map[int32]int)
	matchIndex = make(map[int32]int)

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
			case msg.MsgType == COMMAND:
				r.Inbox() <- msg.Cmd
			case msg.MsgType == REQUESTVOTE:
				if msg.TermId > r.term {
					r.changeTerm(msg.TermId)
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
					r.changeTerm(msg.TermId)
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
						if response.MatchIndex > matchIndex[msg.Pid] {
							matchIndex[msg.Pid] = response.MatchIndex
							if response.MatchIndex == nextIndex[msg.Pid] {
								nextIndex[msg.Pid] = nextIndex[msg.Pid] + 1
							}
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
				
				
				temp :=make([]int,0,1)
				for _,v := range matchIndex {
					temp = append(temp,v)
				}
				sort.Sort(sort.IntSlice(temp))
				if len(temp) % 2 == 0 {
					r.commitIndex = temp[len(temp)/2-1]
				} else {
					r.commitIndex = temp[len(temp)/2]
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
					if r.LastIndex == matchIndex[id] {
						r.TRACE.Println("Sending heartbeat AppendEntry to ", id," ")
						appendEntry = nil
					} else if nextIndex[id] == r.LastIndex+1 {
						r.TRACE.Println("Sending AppendEntry to ", id, " with nil log item at index ", nextIndex[id])
						appendEntry.LogItm = nil
					} else {
						r.TRACE.Println("Sending AppendEntry to ", id, " with log item at index ", nextIndex[id])
						appendEntry.LogItm = r.log[nextIndex[id]]
					}
					msg.LeaderCommit = r.commitIndex
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
						r.TRACE.Println("Sending heartbeat AppendEntry to ", id," ")
						appendEntry = nil
					} else if nextIndex[id] == r.LastIndex+1 {
						r.TRACE.Println("Sending AppendEntry to ", id, " with nil log item at index ", nextIndex[id])
						appendEntry.LogItm = nil
					} else {
						r.TRACE.Println("Sending AppendEntry to ", id, " with log item at index ", nextIndex[id])
						appendEntry.LogItm = r.log[nextIndex[id]]
					}
					msg.LeaderCommit = r.commitIndex
					msg.Msg = appendEntry
					msg.Resp = nil
					r.server.Outbox() <- msg
				}
			}
			r.timer.Reset(time.Duration(2) * time.Second)
		}
	}
}

func (r	*RaftImpl) proc_outbox() {
	for {
		select {
			case logitem := <- r.Outbox():
					if logitem.Data.Type == PUT {
						err := r.db.Put(r.wo,[]byte(logitem.Data.Key), []byte(logitem.Data.Value))
						if err != nil {
							r.TRACE.Println("Put command failed PUT(",logitem.Data.Key,",",logitem.Data.Value,")")
						}
					} else if logitem.Data.Type == DELETE {
						err := r.db.Delete(r.wo, []byte(logitem.Data.Key))
						if err != nil {
							r.TRACE.Println("Delete command failed DELETE(",logitem.Data.Key,")")
						}
					}
		}	
	}
}
