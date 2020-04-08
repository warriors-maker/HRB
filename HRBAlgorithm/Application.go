package HRBAlgorithm

import "time"

func broadcast(byteLength, round int) {
	time.Sleep(3 *time.Second)
	for i := 0; i < round; i++ {
		//if i % 200 == 0 {
		//	time.Sleep(1*time.Second)
		//}
		s := RandStringBytes(byteLength)
		m := MSGStruct{Id: MyID, SenderId:MyID, Data: s, Header:MSG, Round:i}
		for _, server := range serverList {
			//fmt.Println("Protocal send to ", server)
			sendReq := PrepareSend{M: m, SendTo: server}
			SendReqChan <- sendReq
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
