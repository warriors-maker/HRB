package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
)


func startUp() {
	//Set up all the variables needed for the HashReliable Broadcast algorithm
	HRBAlgorithm.AlgorithmSetUp()
	// Read from the file
	netWorkSize := 0
	//list of Faulty nodes the server finds
	//BlackList := make([]string, netWorkSize)
	
}

//Receive a data from the channel that connects to the thread of TcpReader
func filterRecData (ch chan HRBAlgorithm.Message) {
	data := <- ch

	switch v := data.GetHeaderType(); v {
	case HRBAlgorithm.MSG:
		receiveMsg(data)
	case HRBAlgorithm.ECHO:
		receiveEcho(data)
	case HRBAlgorithm.ACC:
		receiveAcc(data)
	case HRBAlgorithm.REQ:
		receiveReq(data)
	case HRBAlgorithm.FWD:
		receiveFwd(data)
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

func cleanUp()  {
	
}


