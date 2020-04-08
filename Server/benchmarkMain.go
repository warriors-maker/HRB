package Server

import (
	"HRB/HRBMessage"
	"time"
)

//used in Benchmark
var internalReadChan chan TcpMessage
var internalWriteChan chan TcpMessage
var externalReadChan chan TcpMessage
var externalWriteChan map[string] messageChan
var statsChan chan HRBMessage.Message
var throughPutBeginTime time.Time



func BenchmarkStart() {
	initChannels()
	initStats()
	statSetup()
	readSetup()
	writeSetup()
}

func initChannels() {
	//Channels between benchmark and Algorithm protocal
	internalReadChan = make (chan TcpMessage, 20000)
	internalWriteChan = make (chan TcpMessage,20000)

	//Channels between benchmark and external nodes in the network
	externalWriteChan = make (map[string] messageChan,20000)
	externalReadChan = make(chan TcpMessage,20000)

	//A Channel for calculating benchmark statistics
	statsChan = make (chan HRBMessage.Message,20000)
	//protocalSendChan = make(chan TcpMessage, 10000)
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

		if data.Message.GetHeaderType() == HRBMessage.Stat {
			statsChan <- data.Message
		} else if data.Message.GetHeaderType() == HRBMessage.MSG {
			if algorithm == 9 {
				if source && sendTo == MyId{
					statsChan <- data.Message
				}
			} else if sendTo == MyId || sendTo == "" || sendTo == "all" {
				statsChan <- data.Message
			}
		}

		// if this is not a Stat Message
		if data.Message.GetHeaderType() != HRBMessage.Stat {
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
	flag := false
	for {
		data := <- externalReadChan
		if ! flag {
			flag = true
			throughPutBeginTime = time.Now()
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

