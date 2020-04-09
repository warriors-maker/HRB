package HRBAlgorithm

import (
	"HRB/HRBMessage"
	"encoding/gob"
	"fmt"
	"time"
)



// Send Phase to the TCPWriter
//var sendChan chan Message

var faulty int
var trusted int
var total int
var MyID string
var algorithm int


var SendReqChan chan HRBMessage.PrepareSend

//key: IP_ID, Value: index in the serverList
var serverMap map[string] int
var serverList []string

//(SenderID + h, bool)
var MessageReceiveSet map[string] bool
//var MessageSentSet map[string] bool

//(SenderId + h, bool)
var EchoReceiveSet map[string] bool
//(sendToID + h, bool)
var EchoSentSet map[string] bool
//Used in Not Simple Version
var EchoRecCountSet map[HRBMessage.ECHOStruct] int


//(SenderId + h, bool)
var AccReceiveSet map[string] bool
var AccSentSet map[string] bool
//(HashStr, list of ids that send Acc)
var AccRecCountSet map[HRBMessage.ACCStruct] []string

var ReqReceiveSet map[string] bool
//(HashStr, list of ids that you send request to)
var ReqSentSet map[HRBMessage.REQStruct] []string


var FwdReceiveSet map[string] bool
//(SendToId, bool)
//var FwdSentSet map[string] bool

//Key: value , Value: Hash(value)
var DataSet map[string] string
var acceptData map[string] bool


/*
Exposed this functions for public used in Server package
*/


func AlgorithmSetUp(myID string, servers []string, trustedCount, faultyCount, round, alg int) {
	round = round / 2
	algorithm = alg


	serverMap = make(map[string] int, round)
	acceptData = make(map[string]bool, round)
	for index, server := range servers {
		serverMap[server] = index
	}
	serverList = servers

	//fmt.Println("These are the servers", serverMap)
	MessageReceiveSet = make(map[string] bool, round)
	//MessageSentSet = make(map[string] bool)

	EchoReceiveSet = make(map[string] bool, round)
	EchoSentSet = make(map[string] bool, round)
	//Used in Acc version
	EchoRecCountSet = make (map[HRBMessage.ECHOStruct] int, round)

	AccReceiveSet = make(map[string] bool, round)
	AccSentSet = make(map[string] bool, round)
	AccRecCountSet = make(map[HRBMessage.ACCStruct] []string, round)

	ReqReceiveSet = make(map[string] bool, round)
	ReqSentSet = make(map[HRBMessage.REQStruct] []string, round)

	FwdReceiveSet = make(map[string] bool, round)
	//FwdSentSet = make(map[string] bool)

	DataSet = make (map[string] string, round)


	//sendChan = make(chan Message)

	SendReqChan = make (chan HRBMessage.PrepareSend)

	//change later based on config
	trusted = trustedCount
	faulty = faultyCount
	total = trusted + faulty
	//fmt.Println("Hey come on:" , trusted, faulty, total)
	MyID = myID

	//augmentRecSend = make(map[string]map[string] [][] digestStruct)
	trustedCount = total


	//Register the concrete type for interface
	gob.Register(HRBMessage.ACCStruct{})
	gob.Register(HRBMessage.FWDStruct{})
	gob.Register(HRBMessage.REQStruct{})
	gob.Register(HRBMessage.MSGStruct{})
	gob.Register(HRBMessage.ECHOStruct{})
	gob.Register(HRBMessage.StatStruct{})
}


func FilterRecData (message HRBMessage.Message) {
	switch v := message.(type) {
	case HRBMessage.MSGStruct:
		Msghandler(message)
	case HRBMessage.ECHOStruct:
		EchoHandler(message)
	case HRBMessage.ACCStruct:
		//fmt.Println("Acc")
		AccHandler(message)
	case HRBMessage.REQStruct:
		//fmt.Println("Req")
		ReqHandler(message)
	case HRBMessage.FWDStruct:
		//fmt.Print("FWD")
		FwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func SimpleBroadcast(byteLength, round int) {
	time.Sleep(3 *time.Second)
	for i := 0; i < round; i++ {
		//if i % 200 == 0 {
		//	time.Sleep(1*time.Second)
		//}
		s := RandStringBytes(byteLength)
		m := HRBMessage.MSGStruct{Id: MyID, SenderId:MyID, Data: s, Header:HRBMessage.MSG, Round:i}
		for _, server := range serverList {
			//fmt.Println("Protocal send to ", server)
			sendReq := HRBMessage.PrepareSend{M: m, SendTo: server}
			SendReqChan <- sendReq
		}
	}
}
