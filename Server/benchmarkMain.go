package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
	"strconv"
	"time"
)

//used in Benchmark
var internalReadChan chan TcpMessage
var internalWriteChan chan TcpMessage
var externalReadChan chan TcpMessage
var externalWriteChan map[string] messageChan
var statsChan chan HRBAlgorithm.Message


func BenchmarkStart() {
	initChannels()
	statSetup()
	readSetup()
	writeSetup()
}

func initChannels() {
	internalReadChan = make (chan TcpMessage)
	internalWriteChan = make (chan TcpMessage)
	externalWriteChan = make (map[string] messageChan)
	externalReadChan = make(chan TcpMessage)
	statsChan = make (chan HRBAlgorithm.Message)
}

func statSetup() {
	go statsCalculate(statsChan)
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
	//count := 0
	go internalReader(internalReadChan)
	for {
		data := <- internalReadChan
		sendTo := data.ID
		//Reliable Broadcast
		if data.Message.GetHeaderType() == HRBAlgorithm.Stat {
			//count += 1
			//fmt.Println(count)
			statsChan <- data.Message
		} else {
			if sendTo == "" || sendTo == "all" {
				for _ , channel := range externalWriteChan {
					//fmt.Println("Send to all now with", id)
					channel <- data
				}
			}  else {
				//fmt.Println("Send to specific now with", sendTo)
				externalWriteChan[sendTo] <- data
			}
		}
	}
}

func networkRead(){
	go ExternalTcpReader(externalReadChan, MyId)
	for {
		data := <- externalReadChan
		if data.Message.GetHeaderType() == HRBAlgorithm.MSG {
			identifier := strconv.Itoa(data.Message.GetRound())
			statsMap[identifier] = time.Now()
		}
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

