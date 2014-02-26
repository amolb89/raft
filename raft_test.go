package main

import (
	clust "github.com/amolb89/raft/clust"
	raft "github.com/amolb89/raft/raftservice"
	"testing"
)

func TestLeaderElection(t *testing.T) {
	configPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/config/serverConfig.json"
	//logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	raft1.Configure(5556, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft2 := new(raft.RaftImpl)
	raft2.Configure(5557, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft3 := new(raft.RaftImpl)
	raft3.Configure(5558, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	raft4 := new(raft.RaftImpl)
	raft4.Configure(5559, configPath, clust.Random(5, 30), 1, CurrLeader, 0.0, 0)

	for {
		select {
		case <-CurrLeader:
			return
		}
	}
}

func TestLeaderDropPackets(t *testing.T) {
	configPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/config/serverConfig.json"
	//logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	CurrLeader := make(chan int, 10)
	raft1 := new(raft.RaftImpl)
	raft1.Configure(5556, configPath, clust.Random(5, 30), 1, CurrLeader, 0.1, 0)

	raft2 := new(raft.RaftImpl)
	raft2.Configure(5557, configPath, clust.Random(5, 30), 1, CurrLeader, 0.2, 0)

	raft3 := new(raft.RaftImpl)
	raft3.Configure(5558, configPath, clust.Random(5, 30), 1, CurrLeader, 0.15, 0)

	raft4 := new(raft.RaftImpl)
	raft4.Configure(5559, configPath, clust.Random(5, 30), 1, CurrLeader, 0.17, 0)

	for {
		select {
		case <-CurrLeader:
			return
		}
	}
}
