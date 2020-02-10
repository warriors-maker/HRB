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

//var SendChans map[string] messageChan
var SendChans chan TcpMessage
var ReadChans chan TcpMessage


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

	time.Sleep(1*time.Second)

	if algorithm == 1 {
		hashSimpleSetup()
		if isSourceFault {

		} else {
			fmt.Println("Running Non-Faulty HashNonEquivocate")
			if source {
				//simpleBroadcast("abcdef")
			}
		}
	} else if algorithm == 2 {
		hashComplexSetup()
		if isSourceFault {
			fmt.Println("Running Faulty HashEquivocate")
			//testSourceFault("abcdef")
		} else {
			fmt.Println("Running Non-Faulty HashEquivocate")
			if source {
				simpleBroadcast("abcdef")
			}
		}
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
}


/*
Seven different algorithms to choose
 */

func hashSimpleSetup() {
	ReadChans := setUpRead()
	go filterSimple(ReadChans)
	go setUpWrite()
}

func hashComplexSetup() {
	ReadChans := setUpRead()
	go filter(ReadChans)
	go setUpWrite()
}


func hashECSimpleSetup() {
	ReadChans := setUpRead()
	go filterSimpleEC(ReadChans)
	go setUpWrite()
}


func hashECComplexSetup() {
	ReadChans := setUpRead()
	go filterComplexEc(ReadChans)
	go setUpWrite()
}

func digestSetup() {
	HRBAlgorithm.InitDigest()
	ReadChans := setUpRead()
	go filterDigest(ReadChans)
	go setUpWrite()
}

func codedSetup() {
	HRBAlgorithm.InitByzCode()
	ReadChans := setUpRead()
	go filterByzCode(ReadChans)
	go setUpWrite()
}

func codedCrashSetup() {
	HRBAlgorithm.InitCrash()
	ReadChans := setUpRead()
	go filterCrashCoded(ReadChans)
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
	ReadChans = make (chan TcpMessage)
	//Start listening data
	go TcpReader(ReadChans, MyId)
	return ReadChans
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
	SendChans = make (messageChan)
	go deliver(MyId, SendChans)
	for {
		req := <- HRBAlgorithm.SendReqChan
		//fmt.Printf("Sending Msg: %+v\n",req.M)
		if req.SendTo == "all" ||  req.SendTo == ""{
			fmt.Println("Sending1")
			tcpMessage := TcpMessage{Message:req.M}
			SendChans <- tcpMessage
		} else {
			fmt.Println("Sending2")
			tcpMessage := TcpMessage{Message:req.M, ID:req.SendTo}
			SendChans <- tcpMessage
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
//				SendChans[server] <- tcpMessage
//			} else if id == 4 || id ==5 {
//
//			} else {
//				codeString := HRBAlgorithm.ConvertBytesToString(shards[id])
//				correct := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, HashData: hashStr, Data: codeString, Header:0, Round:0}
//				tcpMessage := TcpMessage{Message:correct}
//				SendChans[server] <- tcpMessage
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
//				SendChans[server] <- tcpMessage
//			} else if id == 4 || id == 5 {
//
//			} else {
//				tcpMessage := TcpMessage{Message:m}
//				SendChans[server] <- tcpMessage
//			}
//		}
//	}
//}
//

func simpleBroadcast(s string) {
	if source {
		fmt.Println("Broadcast")
		m := HRBAlgorithm.MSGStruct{Id: MyId, SenderId:MyId, Data: s, Header:0, Round:0}
		tcpMessage := TcpMessage{Message:m, ID:"1234"}
		SendChans <- tcpMessage
	}
}



