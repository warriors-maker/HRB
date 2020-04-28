package Server

import (
	"HRB/HRBMessage"
	"fmt"
	"time"
)

//used in Benchmark
var internalReadChan chan TcpMessage
var internalWriteChan chan TcpMessage
var externalReadChan chan TcpMessage
var externalWriteChan map[string] messageChan
var statsChan chan interface{}
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
	statsChan = make (chan interface{},20000)
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

		switch v := data.Message.(type) {
		case HRBMessage.MSGStruct:
			statsChan <- v
			if sendTo == "" || sendTo == "all" {
				fmt.Printf("Sending : %+v\n", v)
				for _ , channel := range externalWriteChan {
					channel <- data
				}
			}  else {
				fmt.Printf("Sending : %+v\n", v)
				externalWriteChan[sendTo] <- data
			}
			break
		case HRBMessage.ECHOStruct:
			if sendTo == "" || sendTo == "all" {
				for _ , channel := range externalWriteChan {
					//fmt.Println("Send to all now with", id)
					channel <- data
				}
			}  else {
				//fmt.Println("Send to specific now with", sendTo)
				externalWriteChan[sendTo] <- data
			}
			break
		case HRBMessage.ACCStruct:
			if sendTo == "" || sendTo == "all" {
				for _ , channel := range externalWriteChan {
					//fmt.Println("Send to all now with", id)
					channel <- data
				}
			}  else {
				//fmt.Println("Send to specific now with", sendTo)
				externalWriteChan[sendTo] <- data
			}
			//fmt.Println("Acc")
			break
		case HRBMessage.REQStruct:
			if sendTo == "" || sendTo == "all" {
				for _ , channel := range externalWriteChan {
					//fmt.Println("Send to all now with", id)
					channel <- data
				}
			}  else {
				//fmt.Println("Send to specific now with", sendTo)
				externalWriteChan[sendTo] <- data
			}
			//fmt.Println("Req")
			break
		case HRBMessage.FWDStruct:
			if sendTo == "" || sendTo == "all" {
				for _ , channel := range externalWriteChan {
					//fmt.Println("Send to all now with", id)
					channel <- data
				}
			}  else {
				//fmt.Println("Send to specific now with", sendTo)
				externalWriteChan[sendTo] <- data
			}
			//fmt.Print("FWD")
			break
		case HRBMessage.StatStruct:
			statsChan <- v
			break;
		default:
			fmt.Printf("Sending : %+v\n", v)
			fmt.Println("I do ot understand what you send")
		}
	}
}

func submitData(val interface{}) {



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

