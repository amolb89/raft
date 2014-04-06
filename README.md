StormDB
====

Distributed and fault-tolerant Key value store

StormDB is a distributed, fault-tolerant key value store, Level DB is used for persistent key value store on single server and the key value store is replicated on multiple server, it is made consistent using Raft consensus algorithm. 

Underlying library is implementation of raft consensus algorithm in golang. It can be used to elect a leader among a group of servers using election algorithm and to replicate leader's log to other servers in a cluster. Log is stored on LevelDb. Clust package is used to send and recieve point to point or broadcast messages in between servers. Messages are transferred using ZeroMQ. All servers starts in follower state. Servers get randomly generated timeout. CurrLeader channel is passed to configure raft, which sends a current leader id on this channel once the leader is elected. Exit channel is used to stop the servers. For testing purpose MockCluster is used to drop a message with certain probability and delay particular message. Logs directory is used to maintain log file(for debugging purpose)

Testing :
Testing package involves running Key value store on 6 servers. Logs are initialized. After the leader is elected, Key- value store on any server is probed and checked to see it doesn't contain any value which is inconsistent. Also new command is sent on any server which redirects it to leader and leader then replicate it to other servers, once it is replicated to majority servers, it is commited. Testing package checks commited record on multiple servers. 

Usage

go build github.com/amolb89/raft/clust

go install github.com/amolb89/raft/clust

go build github.com/amolb89/raft/raftservice

go install github.com/amolb89/raft/raftservice

go test github.com/amolb89/raft

