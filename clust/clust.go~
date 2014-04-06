package clust

import (
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"strconv"
)

const (
	BROADCAST = -1
)

// Message type can be 1-RequestVote 2-Appendentries 3-Ok
type Envelope struct {
	Pid     int32
	MsgType int
	TermId  int32
	LeaderCommit int
	Msg     *AppendEntry
	Resp    *Response
	Cmd 	*Command
}

type Response struct {
	MatchIndex int
	//MatchTerm int
	Success bool
}

type Command struct {
	Type int
	Key string
	Value string
}

type LogItem struct {
	//Index int
	Term int32
	Data *Command
}

type AppendEntry struct {
	PrevLogIndex int
	PrevLogTerm  int32
	LogItm       *LogItem
}

type Server interface {
	Pid() int
	Peers() map[int32]string
	Outbox() chan *Envelope
	Inbox() chan *Envelope
}

type Serv struct {
	pid       int32
	peers     map[int32]string
	outbox    chan *Envelope
	inbox     chan *Envelope
	noOfPeers int
}

func (s *Serv) Set(id int32, configPath string, exit <-chan int) {
	cfig := new(Config)
	//Load configuration from config file
	LoadConfig(configPath, cfig)
	Peers := make(map[int32]string)

	//Initialize list of peers
	s.noOfPeers = 0
	for servers, _ := range cfig.Servers {
		key, _ := strconv.Atoi(servers)
		if int32(key) != id {
			Peers[int32(key)] = cfig.Servers[servers]
			s.noOfPeers++
		} else {
			Peers[int32(key)] = cfig.Servers[servers]
		}
	}

	//Create inbox and outbox channels
	s.outbox = make(chan *Envelope, 100)
	s.inbox = make(chan *Envelope, 100)
	s.pid = id
	s.peers = Peers

	//Start receiving and sending goroutines
	go s.Proc_recv(exit)
	go s.Proc_send(exit)
}

func (s *Serv) NoOfPeers() int {
	return s.noOfPeers
}

func (s *Serv) Pid() int32 {
	return s.pid
}

func (s *Serv) Peers() map[int32]string {
	return s.peers
}

func (s *Serv) Outbox() chan *Envelope {
	return s.outbox
}

func (s *Serv) Inbox() chan *Envelope {
	return s.inbox
}

func (s *Serv) Proc_send(exit <-chan int) {
L:
	for {
		select {
		case <-exit:
			break L
		default:
			msg := <-s.Outbox()
			if msg.Pid == -1 {
				for id, _ := range s.Peers() {
					if id != s.Pid() {
						requester, _ := zmq.NewSocket(zmq.PUSH)
						servAddr := "tcp://" + s.Peers()[id]
						requester.Connect(servAddr)
						msg.Pid = s.Pid()
						msgBytes, _ := json.Marshal(msg)
						requester.SendBytes(msgBytes, 0)
						requester.Close()
					}
				}
			} else {
				requester, _ := zmq.NewSocket(zmq.PUSH)
				servAddr := "tcp://" + s.Peers()[msg.Pid]
				requester.Connect(servAddr)
				msg.Pid = s.Pid()
				msgBytes, _ := json.Marshal(msg)
				requester.SendBytes(msgBytes, 0)
				requester.Close()
			}
		}
	}
}

func (s *Serv) Proc_recv(exit <-chan int) {
	responder, _ := zmq.NewSocket(zmq.PULL)
	defer responder.Close()
	recAddr := "tcp://" + s.Peers()[s.Pid()]
	responder.Bind(recAddr)
L:
	for {
		select {
		case <-exit:
			break L
		default:
			msgBytes, _ := responder.RecvBytes(0)
			var msg Envelope
			json.Unmarshal(msgBytes, &msg)
			s.Inbox() <- &msg
		}
	}
}
