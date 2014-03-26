package clust

import (
	"fmt"
	"time"
)

type MockServer struct {
	dropProb float64
	maxDelay int
	outbox   chan *Envelope
	inbox    chan *Envelope
	server   *Serv
}

func (s *MockServer) Outbox() chan *Envelope {
	return s.outbox
}

func (s *MockServer) Inbox() chan *Envelope {
	return s.inbox
}

func (s *MockServer) NoOfPeers() int {
	return s.server.NoOfPeers()
}

func (s *MockServer) Peers() map[int]string {
	return s.server.Peers()
}

func (s *MockServer) Set(id int, configPath string, DropProb float64, MaxDelay int, exit <-chan int) {
	s.server = new(Serv)
	s.server.Set(id, configPath, exit)
	s.dropProb = DropProb
	s.maxDelay = MaxDelay

	//Create inbox and outbox channels
	s.outbox = make(chan *Envelope, 100)
	s.inbox = make(chan *Envelope, 100)

	go s.Proc_send(exit)
	go s.Proc_recv(exit)
}

func (s *MockServer) Proc_send(exit <-chan int) {
L:
	for {
		select {
		case <-exit:
			break L
		default:
			msg := <-s.Outbox()
			if msg.Pid == -1 {
				for id, _ := range s.server.Peers() {
					if id != s.server.Pid() {
						p := Random(int(s.dropProb*100), 100)
						if p < int(s.dropProb*100) {
							//Drop the packet
							fmt.Println("Outgoing packet dropped at server ", s.server.Pid())
						} else {
							msgNew := new(Envelope)
							msgNew.Pid = id
							msgNew.Msg = msg.Msg
							msgNew.Resp = msg.Resp
							msgNew.TermId = msg.TermId
							msgNew.MsgType = msg.MsgType
							//Delay
							delay := Random(0, s.maxDelay)
							fmt.Println("Outgoing packet delayed by ", delay, " at server ", s.server.Pid())
							time.Sleep(time.Duration(delay) * time.Second)
							s.server.Outbox() <- msgNew
						}
					}
				}
			} else {
				p := Random(int(s.dropProb*100), 100)
				if p < int(s.dropProb*100) {
					//Drop the packet
					fmt.Println("Outgoing packet dropped at server ", s.server.Pid())
				} else {
					delay := Random(0, s.maxDelay)
					fmt.Println("Outgoing packet delayed by ", delay, " at server ", s.server.Pid())
					time.Sleep(time.Duration(delay) * time.Second)
					s.server.Outbox() <- msg
				}
			}
		}
	}
}

func (s *MockServer) Proc_recv(exit <-chan int) {
L:
	for {
		select {
		case <-exit:
			break L
		default:
			msg := <-s.server.Inbox()
			p := Random(int(s.dropProb*100), 100)
			if p < int(s.dropProb*100) {
				//Drop the packet
				fmt.Println("Incoming packet dropped at server ", s.server.Pid())
			} else {
				delay := Random(0, s.maxDelay)
				fmt.Println("Incoming packet delayed by ", delay, " at server ", s.server.Pid())
				time.Sleep(time.Duration(delay) * time.Second)
				s.Inbox() <- msg
			}
		}
	}
}
