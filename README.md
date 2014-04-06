Raft
====

Raft implementation in golang

This library is implementation of raft consensus algorithm in golang. It can be used to elect a leader among a group of servers using election algorithm. It uses a clust package which is used to send and recieve point to point or broadcast messages in between servers. All servers starts in follower state. Servers get randomly generated timeout. CurrLeader channel is passed to configure raft, which sends a current leader id on this channel once the leader is elected. Exit channel is used to stop the servers. For testing purpose MockCluster is used to drop a message with certain probability and delay particular message. Logs directory is used to maintain log file.

Usage

go build github.com/amolb89/raft/clust

go install github.com/amolb89/raft/clust

go build github.com/amolb89/raft/raftservice

go install github.com/amolb89/raft/raftservice

go test github.com/amolb89/raft

