package HRBAlgorithm

import (
	"encoding/gob"
	"fmt"
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

//Key: value , Value: Hash(value)
var DataSet map[string] string

/*
Used in the Erasure Coding
 */
var ecDataSet map[string] [][]byte


// Send Phase to the TCPWriter
//var sendChan chan Message

var faulty int
var trusted int
var total int
var MyID string
var algorithm int

var acceptData map[string] bool

var SendReqChan chan PrepareSend

//key: IP_ID, Value: index in the serverList
var serverMap map[string] int
var serverList []string
var reqSentHash map[string] string
var FwdRecCountSet map[string] int

/*
Digest
 */
var genKey string
//var digestSourceData map[string] string //store the message sent from source
//var digestDataMap map[string] []digestStruct //All the Hash Data from peers

//var digestRecSend map[string] [][]digestStruct //what you have sent in one round
//var faultyCountMap map[string] int
//var faultySet map[string] bool

//var augmentRecSend map[string] map[string] [][]digestStruct //Used during validate step
var binarySet map[string] []Message

var digestTrustCount int

var statsRecord map[string] Stats



func AlgorithmSetUp(myID string, servers []string, trustedCount, faultyCount, round, alg int) {
	round = round / 2
	algorithm = alg
	//fmt.Println(algorithm)
	statsRecord = make(map[string]Stats,round)
	serverMap = make(map[string] int, round)
	acceptData = make(map[string]bool, round)
	FwdRecCountSet = make (map [string] int, round)
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
	EchoRecCountSet = make (map[ECHOStruct] int, round)
	//used in Simple
	simpleEchoRecCountSet = make (map[ECHOStruct] []string, round)

	AccReceiveSet = make(map[string] bool, round)
	AccSentSet = make(map[string] bool, round)
	AccRecCountSet = make(map[ACCStruct] []string, round)

	ReqReceiveSet = make(map[string] bool, round)
	ReqSentSet = make(map[REQStruct] []string, round)

	FwdReceiveSet = make(map[string] bool, round)
	//FwdSentSet = make(map[string] bool)

	DataSet = make (map[string] string, round)

	ecDataSet = make(map[string] [][]byte, round)

	//sendChan = make(chan Message)

	SendReqChan = make (chan PrepareSend)

	//change later based on config
	trusted = trustedCount
	faulty = faultyCount
	total = trusted + faulty
	//fmt.Println("Hey come on:" , trusted, faulty, total)
	MyID = myID
	genKey = MyID
	digestTrustCount = total

	reqSentHash = make(map[string] string, round)
	//digestSourceData = make(map[string] string)
	//digestDataMap = make(map[string] []digestStruct)
	//digestRecSend = make(map[string] [][] digestStruct)
	//faultyCountMap = make(map[string] int)
	//faultySet = make(map[string] bool)

	//notTrustedListMap = make(map[string] []string, round)
	//for _, server := range serverList {
	//	notTrustedListMap[server] = []string{}
	//}


	//augmentRecSend = make(map[string]map[string] [][] digestStruct)
	trustedCount = total
	binarySet = make(map[string] []Message, round)

	//Register the concrete type for interface
	gob.Register(ACCStruct{})
	gob.Register(FWDStruct{})
	gob.Register(REQStruct{})
	gob.Register(MSGStruct{})
	gob.Register(ECHOStruct{})
	gob.Register(Binary{})
	gob.Register(StatStruct{})
}


func SimpleFilterRecData(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		SimpleMsgHandler(message)
	case ECHOStruct:
		SimpleEchoHandler(message)
	case REQStruct:
		SimpleReqHandler(message)
	case FWDStruct:
		SimpleFwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func FilterRecData (message Message) {
	switch v := message.(type) {
	case MSGStruct:
		//fmt.Println("Msg")
		Msghandler(message)
	case ECHOStruct:
		//fmt.Println("Echo")
		EchoHandler(message)
	case ACCStruct:
		//fmt.Println("Acc")
		AccHandler(message)
	case REQStruct:
		//fmt.Println("Req")
		ReqHandler(message)
	case FWDStruct:
		//fmt.Print("FWD")
		FwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func FilterSimpleErasureCodeRecData(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		//fmt.Println("Msg")
		SimpleECMessageHandler(message)
	case ECHOStruct:
		//fmt.Println("Echo")
		SimpleECEchoHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		//fmt.Println("I do ot understand what you send")
	}
}


func FilterComplexErasureRecData(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		//fmt.Println("Msg")
		ComplexECMessageHandler(message)
	case ECHOStruct:
		//fmt.Println("Echo")
		ComplexECEchoHandler(message)
	case ACCStruct:
		//fmt.Println("Acc")
		complexECAccHandler(message)
	case REQStruct:
		//fmt.Println("Req")
		ComplexECReqHandler(message)
	case FWDStruct:
		//fmt.Print("FWD")
		complexECFwdHandler(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I do ot understand what you send")
	}
}

func FilterDigest(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		receivePrepareFromSrc(message)
	case ECHOStruct:
		receiveDigestFromOthers(message)
	case Binary:
		recBinary(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I donot understand what you send")
	}
}

func FilterByzCode(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		ByzRecMsg(message)
	case ECHOStruct:
		ByzRecEcho(message)
	case Binary:
		ByzRecBin(message)
	//case RecSend:
	//	receiveRecSend(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I donot understand what you send")
	}
}

func FilterCrashCode(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		crashRecMsg(message)
	case ECHOStruct:
		crashRecEcho(message)
	case ACCStruct:
		crashRecAcc(message)
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I donot understand what you send")
	}
}

func FilterOptimal(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		if v.GetHeaderType() == MSG_OPT {
			optimalMsgHandler(message)
		} else {
			Msghandler(message)
		}
	case ECHOStruct:
		if v.GetHeaderType() == ECHO_OPT {
			optimalEchoHandler(message)
		} else {
			EchoHandler(message)
		}
	case ACCStruct:
		if v.GetHeaderType() == ACC_OPT {
			optimalAccHandler(message)
		} else {
			AccHandler(message)
		}
	case Binary:
		optimalHashTagHandler(message)
	case FWDStruct:
		if v.GetHeaderType() == FWD_OPT {
			optimalFwdHandler(message)
		} else {
			FwdHandler(message)
		}
	case REQStruct:
		if v.GetHeaderType() == REQ_OPT {
			optimalReqHandler(message)
		} else {
			ReqHandler(message)
		}
	default:
		fmt.Printf("Sending : %+v\n", v)
		fmt.Println("I donot understand what you send")
	}
}

func FilterOptimalAgainst(message Message) {
	switch v := message.(type) {
	case MSGStruct:
		simpleMessageHandler(v)
	}
}

/*
Helper Function
 */

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
		if v == expectedHash {
			return true, k
		}
	}
	return false,""
}