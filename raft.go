package main

import (
	clust "github.com/amolb89/raft/clust"
	raft "github.com/amolb89/raft/raftservice"
	"log"
	"os"
)

var (
	TRACE *log.Logger
	ERROR *log.Logger
)

func main() {
	configPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/config/serverConfig.json"
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	TRACE.Println("Configuring raft servers")

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	raft1.Configure(5556, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft2 := new(raft.RaftImpl)
	raft2.Configure(5557, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft3 := new(raft.RaftImpl)
	raft3.Configure(5558, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft4 := new(raft.RaftImpl)
	raft4.Configure(5559, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

L:
	for {
		select {
		case leader := <-CurrLeader:
			TRACE.Println("Server id ", leader, " is the leader")
			break L
		}
	}
}
