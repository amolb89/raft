package main

import (
	"bufio"
	"log"
	"os"
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
	
}

func raft_normal() {
	
}
