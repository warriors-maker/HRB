package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
	"os"
	"time"
)

/*
Note id is always the Ip + Port for Local
 */

//var protocalSendChan map[string] messageChan
var protocalSendChan chan TcpMessage
var protocalReadChan chan TcpMessage


/*
Only used in Local Mode for debugging purpose
 */

func writeLogFile() {
	file, err := os.OpenFile("output"+MyId+".txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(file, MyId+"\n")

	fmt.Fprint(file, MyId+":6379\n" )
	fmt.Fprint(file, serverList)
	fmt.Fprint(file,"\n")
}

//For One round
func ProtocalStart() {
	fmt.Println("Local Setup")

	//HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
	if isSourceFault {
		trustedCount = trustedCount - 1;
		faultyCount = faultyCount + 1;
		fmt.Println("FaultyCount ", faultyCount)
	}

	//General setup

	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)


	if algorithm == 1 {
		hashSimpleSetup()

	} else if algorithm == 2 {
		hashComplexSetup()

	} else if algorithm == 3 {
		hashECSimpleSetup()
		if isSourceFault {

		} else {
			if source {
				HRBAlgorithm.SimpleECBroadCast("abcdef")
			}
		}
	} else if algorithm == 4 {
		hashECComplexSetup()
		if isSourceFault {
			fmt.Println("Running Faulty ECEquivocate")
			//testEcSourceFault("abcdef")
		} else {
			fmt.Println("Running Non-Faulty ECEquivocate")
			if source {
				HRBAlgorithm.ComplexECBroadCast("abcdef")
			}
		}
	} else if algorithm == 5 {
		digestSetup()
		if source {
			fmt.Println("Digest Broadcast")
			HRBAlgorithm.BroadcastPrepare("abcd", 1)
		}

	} else if algorithm == 6 {
		codedSetup()
		if source {
			fmt.Println("Digest Broadcast")
			HRBAlgorithm.ECByzBroadCast("abcdabcd", 1)
		}

	} else if algorithm == 7 {
		codedCrashSetup()
		if source {
			fmt.Println("Crash Ccoded Broadcast")
			HRBAlgorithm.CrashECBroadCast("abcde", 1)
		}
	} else {
		fmt.Println("Do not understand what you give")
	}

	if algorithm == 1 || algorithm == 2 {
		if source {
			simpleBroadcast("abcdef")
		}
	}
}


/*
Seven different algorithms to choose
 */

func hashSimpleSetup() {
	protocalReadChan = setUpRead()
	go filterSimple(protocalReadChan)
	go setUpWrite()
}

func hashComplexSetup() {
	protocalReadChan = setUpRead()
	go filter(protocalReadChan)
	go setUpWrite()
}


func hashECSimpleSetup() {
	protocalReadChan = setUpRead()
	go filterSimpleEC(protocalReadChan)
	go setUpWrite()
}


func hashECComplexSetup() {
	protocalReadChan = setUpRead()
	go filterComplexEc(protocalReadChan)
	go setUpWrite()
}

func digestSetup() {
	HRBAlgorithm.InitDigest()
	protocalReadChan = setUpRead()
	go filterDigest(protocalReadChan)
	go setUpWrite()
}

func codedSetup() {
	HRBAlgorithm.InitByzCode()
	protocalReadChan= setUpRead()
	go filterByzCode(protocalReadChan)
	go setUpWrite()
}

func codedCrashSetup() {
	HRBAlgorithm.InitCrash()
	protocalReadChan= setUpRead()
	go filterCrashCoded(protocalReadChan)
	go setUpWrite()
}


func NetworkModeStartup(id int) {
	isLocalMode = false
	setUpRead()
	HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
}


/*
Reading from the network
*/

func setUpRead() chan TcpMessage{
	protocalReadChan = make (chan TcpMessage)
	//Start listening data
	go TcpReader(protocalReadChan, MyId)
	return protocalReadChan
}

func filter (ch chan TcpMessage) {
	for {
		message := <-ch
		HRBAlgorithm.FilterRecData(message.Message)
	}
	//Call the filter method
}
func filterSimple(ch chan TcpMessage) {
	for {
		message := <-ch
		HRBAlgorithm.SimpleFilterRecData(message.Message)
	}
	//Call the filter method
}

func filterSimpleEC(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterSimpleErasureCodeRecData(message.Message)
	}
}

func filterComplexEc(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterComplexErasureRecData(message.Message)
	}
}

func filterDigest(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterDigest(message.Message)
	}
}

func filterByzCode(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterByzCode(message.Message)
	}
}

func filterCrashCoded(ch chan TcpMessage) {
	for {
		message := <- ch
		HRBAlgorithm.FilterCrashCode(message.Message)
	}
}

/*
Writing to the Network
*/
//Setting up writting Channels for individual sever



func setUpWrite() {
	protocalSendChan = make(chan TcpMessage)
	go deliver(MyId, protocalSendChan)
	for {
		req := <- HRBAlgorithm.SendReqChan
		//fmt.Printf("Sending Msg: %+v\n",req.M)
		if req.SendTo == "all" ||  req.SendTo == ""{
			fmt.Println("Sending1")
			tcpMessage := TcpMessage{Message:req.M}
			protocalSendChan <- tcpMessage
		} else {
			fmt.Println("Sending2")
			tcpMessage := TcpMessage{Message:req.M, ID:req.SendTo}
			protocalSendChan <- tcpMessage
		}
	}
}



func deliver(ipPort string, ch chan TcpMessage) {
	TcpWriter(ch)
}


/*
Simple Testing
 */

//func testEcSourceFault(s string) {
//	//test
//	if source {
//		shards := HRBAlgorithm.Encode(s, faultyCount + 1, trustedCount - 1)
//		//Get the string version of the string
//		hashStr := HRBAlgorithm.ConvertBytesToString(HRBAlgorithm.Hash([]byte(s)))
//
//		for id , server := range serverList {
//			if id == 2 || id == 3{
//				wrong := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, HashData: hashStr, Data: "", Header:0, Round:0}
//				tcpMessage := TcpMessage{Message: wrong}
//				protocalSendChan[server] <- tcpMessage
//			} else if id == 4 || id ==5 {
//
//			} else {
//				codeString := HRBAlgorithm.ConvertBytesToString(shards[id])
//				correct := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, HashData: hashStr, Data: codeString, Header:0, Round:0}
//				tcpMessage := TcpMessage{Message:correct}
//				protocalSendChan[server] <- tcpMessage
//			}
//		}
//	}
//}
//
//func testSourceFault(s string) {
//	if source {
//		m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data:s, Header:0, Round:0}
//		faultym :=  HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data:"", Header:0, Round:0}
//		for id , server := range serverList {
//			if id == 2 || id == 3 {
//				tcpMessage := TcpMessage{Message:faultym}
//				protocalSendChan[server] <- tcpMessage
//			} else if id == 4 || id == 5 {
//
//			} else {
//				tcpMessage := TcpMessage{Message:m}
//				protocalSendChan[server] <- tcpMessage
//			}
//		}
//	}
//}
//

func simpleBroadcast(s string) {
	if source {
		time.Sleep(1*time.Second)
		fmt.Println("Broadcast")
		m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data: s, Header:0, Round:0}
		tcpMessage := TcpMessage{Message:m}
		protocalSendChan <- tcpMessage
	}
}



