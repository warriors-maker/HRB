package main

import (
	"HRB/HRBAlgorithm"
	"HRB/HRBMessage"
	"time"
)

func broadcast(byteLength, round int) {
	time.Sleep(3 *time.Second)
	for i := 0; i < round; i++ {
		//if i % 200 == 0 {
		//	time.Sleep(1*time.Second)
		//}
		s := HRBAlgorithm.RandStringBytes(byteLength)
		m := HRBMessage.MSGStruct{Id: HRBAlgorithm.MyID, SenderId: HRBAlgorithm.MyID, Data: s, Header: HRBMessage.MSG, Round:i}
		for _, server := range HRBAlgorithm.serverList {
			//fmt.Println("Protocal send to ", server)
			sendReq := HRBMessage.PrepareSend{M: m, SendTo: server}
			HRBAlgorithm.SendReqChan <- sendReq
		}
	}
}

func generateData() {

}

func readFromAlg() {

}

func sendToAlg() {

}

func sendStatsToBenchmark() {

}
