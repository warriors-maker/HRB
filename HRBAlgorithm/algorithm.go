package HRBAlgorithm

import (
	"encoding/gob"
	"fmt"
	"sync"
)

//(SenderID + h, bool)
var MessageReceiveSet map[string] bool
//var MessageSentSet map[string] bool

//(SenderId + h, bool)
var EchoReceiveSet map[string] bool
//(sendToID + h, bool)
var EchoSentSet map[string] bool
//Used in Not Simple Version
var EchoRecCountSet map[ECHOStruct] int
//Used in SimpleVersion because of the place where they get accept
var simpleEchoRecCountSet map[ECHOStruct] []string

//(SenderId + h, bool)
var AccReceiveSet map[string] bool
var AccSentSet map[string] bool
//(HashStr, list of ids that send Acc)
var AccRecCountSet map[ACCStruct] []string

var ReqReceiveSet map[string] bool
//(HashStr, list of ids that you send request to)
var ReqSentSet map[REQStruct] []string


var FwdReceiveSet map[string] bool
//(SendToId, bool)
//var FwdSentSet map[string] bool

var DataSet map[string] *receiveData
// Send Phase to the TCPWriter
var sendChan chan Message

var faulty int
var trusted int
var total int
var MyID string

var acceptData []interface{}

var SendReqChan chan PrepareSend


type receiveData struct {
	mux sync.Mutex
	isDirect bool
	hashValue string
	count int
}

func (d *receiveData) update(message Message, hashString string) {
	fmt.Printf("Update data: %+v\n", message)
	d.mux.Lock()
	if message.GetHeaderType() == MSG {
		d.isDirect = true
		d.hashValue = hashString
	} else if message.GetHeaderType() == FWD {
		if ! d.isDirect {
			d.count = d.count + 1
			d.hashValue = hashString
		}
	}
	d.mux.Unlock()
}

func (d *receiveData) fetchInfo() (bool, string, int){
	return d.isDirect, d.hashValue, d.count
}



func AlgorithmSetUp(myID string, serverList []string, trustedCount, faultyCount int) {
	MessageReceiveSet = make(map[string] bool)
	//MessageSentSet = make(map[string] bool)

	EchoReceiveSet = make(map[string] bool)
	EchoSentSet = make(map[string] bool)
	EchoRecCountSet = make (map[ECHOStruct] int)
	//used in Simple
	simpleEchoRecCountSet = make (map[ECHOStruct] []string)

	AccReceiveSet = make(map[string] bool)
	AccSentSet = make(map[string] bool)
	AccRecCountSet = make(map[ACCStruct] []string)

	ReqReceiveSet = make(map[string] bool)
	ReqSentSet = make(map[REQStruct] []string)

	FwdReceiveSet = make(map[string] bool)
	//FwdSentSet = make(map[string] bool)

	DataSet = make (map[string] *receiveData)

	sendChan = make(chan Message)

	SendReqChan = make (chan PrepareSend)

	//change later based on config
	trusted = trustedCount
	faulty = faultyCount
	total = trusted + faulty
	MyID = myID

	//Register the concrete type for interface
	gob.Register(ACCStruct{})
	gob.Register(FWDStruct{})
	gob.Register(REQStruct{})
	gob.Register(MSGStruct{})
	gob.Register(ECHOStruct{})
}

func checkRecMsg(id string) bool{
	for k,_ := range MessageReceiveSet {
		if k == id {
			return true
		}
	}
	return false
}

func hasSent(l []string, val string) bool{
	for _, v := range l {
		if v == val {
			return true
		}
	}
	return false
}

func checkDataExist(expectedHash string) (bool, string) {
	for k,v := range DataSet {
		if v.hashValue == expectedHash {
			//fmt.Println("Check exist" + expectedHash)
			return true, k
		}
	}
	return false,""
}








