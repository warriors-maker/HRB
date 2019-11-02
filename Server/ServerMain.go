package Server

import (
	"HRB/HRBAlgorithm"
	"HRB/Network"
	"fmt"
)
type messageChan chan Network.TcpMessage

var serverList []string
var MyId string //basically the IP of individual server
var isFault bool //check whether I should be the faulty based on configuration

var sendChans []messageChan
var readChans chan Network.TcpMessage

var isLocalMode bool //indicate whether this is a local mode
var source bool //A flag to indicate whether I am the sender

/*
Only used in Local Mode
 */
var localId int

func LocalModeStartup(id int) {
	isLocalMode = true
	localId = id

	peerStartup(isLocalMode)
	if serverList[0] == MyId {
		source = true
	} else {
		source = false
	}
	setUpRead()
	setUpWrite()
}

func NetworkModeStartup() {
	isLocalMode = false
	peerStartup(isLocalMode)

	setUpRead()
	setUpWrite()
}

//Start up the peer
func peerStartup(local bool) {
	trustedPath := "./Configuration/trusted"
	faultedPath := "./Configuration/faulty"
	if local {
		serverList, MyId, isFault = readServerListLocal(trustedPath, faultedPath, localId)
	} else {
		serverList, MyId, isFault = readServerListNetwork(trustedPath, faultedPath)
	}
}


/*
Reading from the network
 */

func setUpRead() {
	readChans = make (chan Network.TcpMessage)
	//Start listening data
	go Network.TcpReader(readChans, MyId)
	//Channel that filters the data based on the message type
	go filterRecData(readChans)

}

func filterRecData (ch chan Network.TcpMessage) {
	for {
		message := <- ch
		data := message.Message

		switch v := data.(type) {
		case HRBAlgorithm.MSGStruct:
			fmt.Println("Msg")
			//receiveMsg(data)
		case HRBAlgorithm.ECHOStruct:
			fmt.Println("Echo")
			//receiveEcho(data)
		case HRBAlgorithm.ACCStruct:
			fmt.Println("Acc")
			//receiveAcc(data)
		case HRBAlgorithm.REQStruct:
			fmt.Println("Req")
			//receiveReq(data)
		case HRBAlgorithm.FWDStruct:
			fmt.Print("FWD")
		default:
			fmt.Printf("Sending : %+v\n", v)
			fmt.Println("I do ot understand what you send")
		}
	}
}

func receiveMsg (data HRBAlgorithm.Message) {
	HRBAlgorithm.Msghandler(data)
}

func receiveEcho (data HRBAlgorithm.Message) {
	HRBAlgorithm.EchoHandler(data)
}

func receiveAcc (data HRBAlgorithm.Message) {
	HRBAlgorithm.AccHandler(data)
}

func receiveReq(data HRBAlgorithm.Message) {
	HRBAlgorithm.ReqHandler(data)
}

func receiveFwd (data HRBAlgorithm.Message) {
	HRBAlgorithm.FwdHandler(data)
}

/*
Writing to the Network
*/
//Setting up writting Channels for individual sever

func setUpWrite() {
	sendChans = make ([]messageChan, len(serverList))

	for i := range sendChans {
		sendChans[i] = make(chan Network.TcpMessage)
		//First Param: the server you are writting to
		//Sending channel
		go deliver(serverList[i], sendChans[i])
	}

}

func deliver(ipPort string, ch chan Network.TcpMessage) {
	Network.TcpWriter(ipPort, ch)
}

func simpleTest() {
	//test
	if source {
		var m HRBAlgorithm.Message

		m = HRBAlgorithm.FWDStruct{Id:MyId, SenderId:MyId}
		message := Network.TcpMessage{Message: m, ID:MyId}
		for _,ch := range sendChans {
			ch <- message
		}

		m = HRBAlgorithm.ACCStruct{Id:MyId, SenderId:MyId}
		message = Network.TcpMessage{Message: m, ID:MyId}
		for _,ch := range sendChans {
			ch <- message
		}

		m = HRBAlgorithm.REQStruct{Id:MyId, SenderId:MyId}
		message = Network.TcpMessage{Message: m, ID:MyId}

		for _,ch := range sendChans {
			ch <- message
		}

		m = HRBAlgorithm.ECHOStruct{Id:MyId, SenderId:MyId}
		message = Network.TcpMessage{Message: m, ID:MyId}

		for _,ch := range sendChans {
			ch <- message
		}

		m = HRBAlgorithm.MSGStruct{Id:MyId, SenderId:MyId}
		message = Network.TcpMessage{Message: m, ID:MyId}
		for _,ch := range sendChans {
			ch <- message
		}
	}
}


