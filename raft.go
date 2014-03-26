package main

import (
	"bufio"
	"encoding/json"
	clust "github.com/amolb89/raft/clust"
	raft "github.com/amolb89/raft/raftservice"
	"io/ioutil"
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
	exit1 := make(chan int, 1)
	l1 := new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 := new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 := new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 := new(clust.LogItem)
	l4.Term = 4
	l4.Data = "move4"
	l5 := new(clust.LogItem)
	l5.Term = 4
	l5.Data = "move5"
	l6 := new(clust.LogItem)
	l6.Term = 5
	l6.Data = "move6"
	l7 := new(clust.LogItem)
	l7.Term = 5
	l7.Data = "move7"
	l8 := new(clust.LogItem)
	l8.Term = 6
	l8.Data = "move8"
	l9 := new(clust.LogItem)
	l9.Term = 6
	l9.Data = "move9"
	c1 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4, l5, l6, l7, l8, l9}}
	logBytes1, _ := json.Marshal(c1)
	logFile1 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5556"
	err1 := ioutil.WriteFile(logFile1, logBytes1, 0644)
	if err1 != nil {
		panic(err1)
	}
	raft1.Configure(5556, configPath, clust.Random(5, 30), 6, CurrLeader, 0.0, 0, exit1)

	raft2 := new(raft.RaftImpl)
	exit2 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 = new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 = new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 = new(clust.LogItem)
	l4.Term = 4
	l4.Data = "move4"
	c2 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4}}
	logBytes2, _ := json.Marshal(c2)
	logFile2 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5557"
	err2 := ioutil.WriteFile(logFile2, logBytes2, 0644)
	if err2 != nil {
		panic(err2)
	}
	raft2.Configure(5557, configPath, clust.Random(5, 30), 4, CurrLeader, 0.0, 0, exit2)

	raft3 := new(raft.RaftImpl)
	exit3 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 = new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 = new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 = new(clust.LogItem)
	l4.Term = 4
	l4.Data = "move4"
	l5 = new(clust.LogItem)
	l5.Term = 4
	l5.Data = "move5"
	l6 = new(clust.LogItem)
	l6.Term = 5
	l6.Data = "move6"
	l7 = new(clust.LogItem)
	l7.Term = 5
	l7.Data = "move7"
	l8 = new(clust.LogItem)
	l8.Term = 6
	l8.Data = "move8"
	l9 = new(clust.LogItem)
	l9.Term = 6
	l9.Data = "move9"
	l10 := new(clust.LogItem)
	l10.Term = 6
	l10.Data = "move10"
	l11 := new(clust.LogItem)
	l11.Term = 6
	l11.Data = "move123"
	c3 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11}}
	logBytes3, _ := json.Marshal(c3)
	logFile3 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5558"
	err3 := ioutil.WriteFile(logFile3, logBytes3, 0644)
	if err3 != nil {
		panic(err3)
	}
	raft3.Configure(5558, configPath, clust.Random(5, 30), 6, CurrLeader, 0.0, 0, exit3)

	raft4 := new(raft.RaftImpl)
	exit4 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 = new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 = new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 = new(clust.LogItem)
	l4.Term = 4
	l4.Data = "move4"
	l5 = new(clust.LogItem)
	l5.Term = 4
	l5.Data = "move5"
	l6 = new(clust.LogItem)
	l6.Term = 5
	l6.Data = "move6"
	l7 = new(clust.LogItem)
	l7.Term = 5
	l7.Data = "move7"
	l8 = new(clust.LogItem)
	l8.Term = 6
	l8.Data = "move8"
	l9 = new(clust.LogItem)
	l9.Term = 6
	l9.Data = "move9"
	l10 = new(clust.LogItem)
	l10.Term = 6
	l10.Data = "mov9"
	l11 = new(clust.LogItem)
	l11.Term = 7
	l11.Data = "mo9"
	l12 := new(clust.LogItem)
	l12.Term = 7
	l12.Data = "mo19"
	c4 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11, l12}}
	logBytes4, _ := json.Marshal(c4)
	logFile4 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5559"
	err4 := ioutil.WriteFile(logFile4, logBytes4, 0644)
	if err4 != nil {
		panic(err4)
	}
	raft4.Configure(5559, configPath, clust.Random(5, 30), 7, CurrLeader, 0.0, 0, exit4)

	raft5 := new(raft.RaftImpl)
	exit5 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 = new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 = new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 = new(clust.LogItem)
	l4.Term = 4
	l4.Data = "move4"
	l5 = new(clust.LogItem)
	l5.Term = 4
	l5.Data = "move5"
	l6 = new(clust.LogItem)
	l6.Term = 4
	l6.Data = "move51"
	l7 = new(clust.LogItem)
	l7.Term = 4
	l7.Data = "move52"
	c5 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4, l5, l6, l7}}
	logBytes5, _ := json.Marshal(c5)
	logFile5 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5560"
	err5 := ioutil.WriteFile(logFile5, logBytes5, 0644)
	if err5 != nil {
		panic(err5)
	}
	raft5.Configure(5560, configPath, clust.Random(5, 30), 4, CurrLeader, 0.0, 0, exit5)

	raft6 := new(raft.RaftImpl)
	exit6 := make(chan int, 1)
	l1 = new(clust.LogItem)
	l1.Term = 1
	l1.Data = "cp"
	l2 = new(clust.LogItem)
	l2.Term = 1
	l2.Data = "move"
	l3 = new(clust.LogItem)
	l3.Term = 1
	l3.Data = "mov"
	l4 = new(clust.LogItem)
	l4.Term = 2
	l4.Data = "move4"
	l5 = new(clust.LogItem)
	l5.Term = 2
	l5.Data = "move5"
	l6 = new(clust.LogItem)
	l6.Term = 2
	l6.Data = "move51"
	l7 = new(clust.LogItem)
	l7.Term = 3
	l7.Data = "move52"
	l8 = new(clust.LogItem)
	l8.Term = 3
	l8.Data = "move8"
	l9 = new(clust.LogItem)
	l9.Term = 3
	l9.Data = "move9"
	l10 = new(clust.LogItem)
	l10.Term = 3
	l10.Data = "move10"
	l11 = new(clust.LogItem)
	l11.Term = 3
	l11.Data = "move123"
	c6 := raft.Log{[]*clust.LogItem{l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11}}
	logBytes6, _ := json.Marshal(c6)
	logFile6 := "/home/amol/Desktop/raft/src/github.com/amolb89/raft/log@5561"
	err6 := ioutil.WriteFile(logFile6, logBytes6, 0644)
	if err6 != nil {
		panic(err6)
	}
	raft6.Configure(5561, configPath, clust.Random(5, 30), 3, CurrLeader, 0.0, 0, exit6)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
