package main

import (
	"bufio"
	"encoding/json"
	clust "github.com/amolb89/raft/clust"
	raft "github.com/amolb89/raft/raftservice"
	levigo "github.com/jmhodges/levigo"
	"log"
	"os"
	"time"
)

const (
	GET     = 1
	PUT     = 2
	DELETE  = 3
	COMMAND = 5
)

var (
	TRACE *log.Logger
	ERROR *log.Logger
)

func main() {
	//raft_init()
	raft_normal()
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func raft_init() {
	configPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/config/serverConfig.json"
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3 << 30))
	opts.SetCreateIfMissing(true)
	wo := levigo.NewWriteOptions()

	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	TRACE.Println("Configuring raft servers")
	CurrLeader := make(chan int32, 10)

	kvstore1 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5556"
	db, _ := levigo.Open(kvstore1, opts)
	raft1 := new(raft.RaftImpl)
	exit1 := make(chan int, 1)
	l1 := new(clust.LogItem)
	l1.Term = 1
	cmd := new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	var index int32
	index = 1
	logItemBytes, _ := json.Marshal(l1)
	indexBytes, _ := json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l2 := new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l3 := new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l4 := new(clust.LogItem)
	l4.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH09"
	cmd.Value = "Kolhapur"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l5 := new(clust.LogItem)
	l5.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH10"
	cmd.Value = "Sangli"
	l5.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l5)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l6 := new(clust.LogItem)
	l6.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH13"
	cmd.Value = "Solapur"
	l6.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l6)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l7 := new(clust.LogItem)
	l7.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH18"
	cmd.Value = "Dhule"
	l7.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l7)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l8 := new(clust.LogItem)
	l8.Term = 6
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH24"
	cmd.Value = "Latur"
	l8.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l8)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)

	l9 := new(clust.LogItem)
	l9.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH18"
	l9.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l9)
	indexBytes, _ = json.Marshal(index)
	db.Put(wo, indexBytes, logItemBytes)
	db.Close()
	raft1.Configure(5556, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit1)

	kvstore1 = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5557"
	db1, _ := levigo.Open(kvstore1, opts)
	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	index = 1
	logItemBytes, _ = json.Marshal(l1)
	indexBytes, _ = json.Marshal(index)
	db1.Put(wo, indexBytes, logItemBytes)

	l2 = new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db1.Put(wo, indexBytes, logItemBytes)

	l3 = new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db1.Put(wo, indexBytes, logItemBytes)

	l4 = new(clust.LogItem)
	l4.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH09"
	cmd.Value = "Kolhapur"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db1.Put(wo, indexBytes, logItemBytes)
	db1.Close()
	raft2.Configure(5557, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit2)

	kvstore1 = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5558"
	db3, _ := levigo.Open(kvstore1, opts)
	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	index = 1
	logItemBytes, _ = json.Marshal(l1)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l2 = new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l3 = new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l4 = new(clust.LogItem)
	l4.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH09"
	cmd.Value = "Kolhapur"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l5 = new(clust.LogItem)
	l5.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH10"
	cmd.Value = "Sangli"
	l5.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l5)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l6 = new(clust.LogItem)
	l6.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH13"
	cmd.Value = "Solapur"
	l6.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l6)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l7 = new(clust.LogItem)
	l7.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH18"
	cmd.Value = "Dhule"
	l7.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l7)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l8 = new(clust.LogItem)
	l8.Term = 6
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH24"
	cmd.Value = "Latur"
	l8.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l8)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l9 = new(clust.LogItem)
	l9.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH18"
	l9.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l9)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l10 := new(clust.LogItem)
	l10.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH24"
	l10.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l10)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	l11 := new(clust.LogItem)
	l11.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH05"
	l11.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l11)
	indexBytes, _ = json.Marshal(index)
	db3.Put(wo, indexBytes, logItemBytes)

	db3.Close()
	raft3.Configure(5558, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit3)

	kvstore1 = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5559"
	db4, _ := levigo.Open(kvstore1, opts)
	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	index = 1
	logItemBytes, _ = json.Marshal(l1)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l2 = new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l3 = new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l4 = new(clust.LogItem)
	l4.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH09"
	cmd.Value = "Kolhapur"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l5 = new(clust.LogItem)
	l5.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH10"
	cmd.Value = "Sangli"
	l5.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l5)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l6 = new(clust.LogItem)
	l6.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH13"
	cmd.Value = "Solapur"
	l6.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l6)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l7 = new(clust.LogItem)
	l7.Term = 5
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH18"
	cmd.Value = "Dhule"
	l7.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l7)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l8 = new(clust.LogItem)
	l8.Term = 6
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH24"
	cmd.Value = "Latur"
	l8.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l8)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l9 = new(clust.LogItem)
	l9.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH18"
	l9.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l9)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l10 = new(clust.LogItem)
	l10.Term = 6
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH24"
	l10.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l10)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l11 = new(clust.LogItem)
	l11.Term = 7
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH32"
	cmd.Value = "Wardha"
	l11.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l11)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)

	l12 := new(clust.LogItem)
	l12.Term = 7
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH35"
	cmd.Value = "Gondia"
	l12.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l12)
	indexBytes, _ = json.Marshal(index)
	db4.Put(wo, indexBytes, logItemBytes)
	db4.Close()
	raft4.Configure(5559, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit4)

	kvstore1 = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5560"
	db5, _ := levigo.Open(kvstore1, opts)
	raft5 := new(raft.RaftImpl)
	exit5 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	index = 1
	logItemBytes, _ = json.Marshal(l1)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l2 = new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l3 = new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l4 = new(clust.LogItem)
	l4.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH09"
	cmd.Value = "Kolhapur"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l5 = new(clust.LogItem)
	l5.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH10"
	cmd.Value = "Sangli"
	l5.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l5)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l6 = new(clust.LogItem)
	l6.Term = 4
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH21"
	cmd.Value = "Jalna"
	l6.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l6)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)

	l7 = new(clust.LogItem)
	l7.Term = 4
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH12"
	l7.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l7)
	indexBytes, _ = json.Marshal(index)
	db5.Put(wo, indexBytes, logItemBytes)
	db5.Close()
	raft5.Configure(5560, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit5)

	kvstore1 = "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5561"
	db6, _ := levigo.Open(kvstore1, opts)
	raft6 := new(raft.RaftImpl)
	exit6 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH12"
	cmd.Value = "Pune"
	l1.Data = cmd
	index = 1
	logItemBytes, _ = json.Marshal(l1)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l2 = new(clust.LogItem)
	l2.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH01"
	cmd.Value = "South Mumbai"
	l2.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l2)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l3 = new(clust.LogItem)
	l3.Term = 1
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH05"
	cmd.Value = "Kalyan"
	l3.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l3)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l4 = new(clust.LogItem)
	l4.Term = 2
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MN09"
	cmd.Value = "Imphal"
	l4.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l4)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l5 = new(clust.LogItem)
	l5.Term = 2
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "ML05"
	cmd.Value = "Shillong"
	l5.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l5)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l6 = new(clust.LogItem)
	l6.Term = 2
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH44"
	cmd.Value = "Beed"
	l6.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l6)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l7 = new(clust.LogItem)
	l7.Term = 3
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH48"
	cmd.Value = "Vasai"
	l7.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l7)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l8 = new(clust.LogItem)
	l8.Term = 3
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH51"
	cmd.Value = "Nashik"
	l8.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l8)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l9 = new(clust.LogItem)
	l9.Term = 3
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH43"
	cmd.Value = "Vashi"
	l9.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l9)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l10 = new(clust.LogItem)
	l10.Term = 3
	cmd = new(clust.Command)
	cmd.Type = DELETE
	cmd.Key = "MH05"
	l10.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l10)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)

	l11 = new(clust.LogItem)
	l11.Term = 3
	cmd = new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "MH31"
	cmd.Value = "Nagpur"
	l11.Data = cmd
	index = index + 1
	logItemBytes, _ = json.Marshal(l11)
	indexBytes, _ = json.Marshal(index)
	db6.Put(wo, indexBytes, logItemBytes)
	db6.Close()
	raft6.Configure(5561, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit6)
}

func raft_normal() {
	configPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/config/serverConfig.json"
	logPath := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/logs/log"

	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	TRACE = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	TRACE.Println("Configuring raft servers")
	CurrLeader := make(chan int32, 10)
	raft1 := new(raft.RaftImpl)
	exit1 := make(chan int, 1)

	raft1.Configure(5556, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit1)

	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	raft2.Configure(5557, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit2)

	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	raft3.Configure(5558, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit3)

	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	raft4.Configure(5559, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit4)

	raft5 := new(raft.RaftImpl)
	exit5 := make(chan int, 1)
	raft5.Configure(5560, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit5)

	raft6 := new(raft.RaftImpl)
	exit6 := make(chan int, 1)
	raft6.Configure(5561, configPath, clust.Random(5, 30), CurrLeader, 0.0, 0, exit6)

	time.Sleep(32 * time.Second)
	msg := new(clust.Envelope)
	msg.MsgType = COMMAND
	cmd := new(clust.Command)
	cmd.Type = PUT
	cmd.Key = "DL01"
	cmd.Value = "New Delhi"
	msg.Cmd = cmd
	raft6.Server().Inbox() <- msg
}
