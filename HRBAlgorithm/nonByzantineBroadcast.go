package HRBAlgorithm

import (
	"time"
)



func NonFaultyBroadCast(length , round int) {
	//need to make sure that coded element > f
	time.Sleep(3*time.Second)

	for r := 0; r < round; r ++ {
		s := RandStringBytes(length)

		for i := 0; i < total; i++ {
			//fmt.Println("Send")
			m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: s, Round: r}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}

	}
}

func nonFaultyMessageHandler(message Message) {
	//fmt.Println("Receive")
	stats := StatStruct{Id:message.GetId(), Round: message.GetRound(), Header:Stat}
	statInfo :=PrepareSend{M:stats, SendTo:MyID}
	SendReqChan <- statInfo
}
