package main

import (
	clust "github.com/amolb89/raft/clust"
	raft "github.com/amolb89/raft/raftservice"
	"testing"
)

func TestLeaderDropPackets(t *testing.T) {
	configPath := "./config/serverConfig.json"
	//logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	exit1 := make(chan int, 1)
	raft1.Configure(5556, configPath, clust.Random(5, 30), 1, CurrLeader, 0.1, 0, exit1)

	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	raft2.Configure(5557, configPath, clust.Random(5, 30), 1, CurrLeader, 0.2, 0, exit2)

	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	raft3.Configure(5558, configPath, clust.Random(5, 30), 1, CurrLeader, 0.15, 0, exit3)

	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	raft4.Configure(5559, configPath, clust.Random(5, 30), 1, CurrLeader, 0.17, 0, exit4)

	<-CurrLeader
	exit1 <- 1
	exit2 <- 1
	exit3 <- 1
	exit4 <- 1
}

func TestLeaderElection(t *testing.T) {
	configPath := "./config/serverConfig2.json"
	//logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	exit1 := make(chan int, 1)
	raft1.Configure(5565, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit1)

	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	raft2.Configure(5566, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit2)

	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	raft3.Configure(5567, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit3)

	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	raft4.Configure(5568, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit4)

	<-CurrLeader
	exit1 <- 1
	exit2 <- 1
	exit3 <- 1
	exit4 <- 1
}

func TestFailedServer(t *testing.T) {
	configPath := "./config/serverConfig1.json"
	//logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	exit1 := make(chan int, 1)
	raft1.Configure(5561, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit1)

	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	raft2.Configure(5562, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit2)

	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	raft3.Configure(5563, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0, exit3)

	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	raft4.Configure(5564, configPath, clust.Random(5, 30), 1, CurrLeader, 1.0, 0, exit4)

	<-CurrLeader
	exit1 <- 1
	exit2 <- 1
	exit3 <- 1
	exit4 <- 1
}
