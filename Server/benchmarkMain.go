package Server

import "fmt"

//used in Benchmark
var internalReadChan chan TcpMessage
var internalWriteChan chan TcpMessage
var externalReadChan chan TcpMessage
var externalWriteChan map[string] messageChan

func BenchmarkStart() {
	initChannels()
	readSetup()
	writeSetup()
}

func initChannels() {
	internalReadChan = make (chan TcpMessage)
	internalWriteChan = make (chan TcpMessage)
	externalWriteChan = make (map[string] messageChan)
	externalReadChan = make(chan TcpMessage)
}

func readSetup() {
	go internalRead()
	go networkRead()
}

func writeSetup() {
	go internalWrite()
	go networkWrite()
}


func internalRead() {
	go internalReader(internalReadChan)
	for {
		data := <- internalReadChan
		sendTo := data.ID
		if sendTo == "" || sendTo == "all" {
			for id , channel := range externalWriteChan {
				fmt.Println("Send to all now with", id)
				channel <- data
			}
		}  else {
			fmt.Println("Send to specific now with", sendTo)
			externalWriteChan[sendTo] <- data
		}
	}
}

func networkRead(){
	go ExternalTcpReader(externalReadChan, MyId)
	for {
		data := <- externalReadChan
		internalWriteChan <- data
	}
}

func internalWrite() {
	internalWriter(protocalReadAddr, internalWriteChan)
}

func networkWrite() {
	//Responsible for writing to other servers
	for _, serverId := range serverList {
		externalWriteChan[serverId] = make(chan TcpMessage)
		go BenchmarkDeliver(serverId, externalWriteChan[serverId])
	}
}


func BenchmarkDeliver(ipPort string, ch chan TcpMessage) {
	ExternalTcpWriter(ipPort, ch)
}

