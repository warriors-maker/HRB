package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
	"math/rand"
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


//For One round
func ProtocalStart() {
	//fmt.Println("Local Setup")

	//HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount)
	if isSourceFault {
		trustedCount = trustedCount - 1;
		faultyCount = faultyCount + 1;
		//fmt.Println("FaultyCount ", faultyCount)
		if MyId == serverList[0] {
			isFault = true
			isSourceFault = true
		}
	}

	protocalReadChan = setUpRead()
	if algorithm == 1 {
		//Initialize the variable for your algorithm
		HRBAlgorithm.AlgorithmSetUp(MyId, serverList, trustedCount, faultyCount, round, algorithm)
		//set up the read write channel
		go filter(protocalReadChan)
		if source {
			go HRBAlgorithm.SimpleBroadcast(dataSize, round)
		}
	} else {
		fmt.Println("Do not understand what you give")
	}

	go setUpWrite()

}

func filter(ch chan TcpMessage) {
	for {
		message := <-ch
		HRBAlgorithm.FilterRecData(message.Message)
	}
}

/*
Reading from the network
*/

func setUpRead() chan TcpMessage{
	protocalReadChan = make (chan TcpMessage,20000)
	//Start listening data
	go TcpReader(protocalReadChan, MyId)
	return protocalReadChan
}


/*
Writing to the Network
*/
//Setting up writting Channels for individual sever



func setUpWrite() {
	protocalSendChan = make(chan TcpMessage,20000)
	go TcpWriter(protocalSendChan)

	for {
		req := <- HRBAlgorithm.SendReqChan
		if req.SendTo == ""{
			for _, id := range serverList {
				tcpMessage := TcpMessage{Message:req.M, ID:id}
				protocalSendChan <- tcpMessage
			}
		} else {
			//fmt.Println("Protocal send to ", req.SendTo)
			tcpMessage := TcpMessage{Message:req.M, ID:req.SendTo}
			protocalSendChan <- tcpMessage
		}
	}
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}







