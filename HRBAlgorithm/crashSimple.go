package HRBAlgorithm

import (
	"time"
)

var simpleMsgRec map[string] bool

func InitSimpleCrash() {
	simpleMsgRec = make(map[string] bool)
}

func CrashBroadCast(length , round int) {
	//need to make sure that coded element > f
	time.Sleep(3*time.Second)

	for r := 0; r < round; r ++ {
		s := RandStringBytes(length)

		for i := 0; i < total; i++ {
			m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: s, Round: r}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}

	}
}

func simpleMessageHandler(message Message) {
	identifier := identifierCreate(message.GetId(), message.GetRound())
	if MyID == serverList[0]{
		if _, e := simpleMsgRec[identifier]; !e {
			simpleMsgRec[identifier] = true
			stats := StatStruct{Id:message.GetId(), Round: message.GetRound(), Header:Stat}
			statInfo :=PrepareSend{M:stats, SendTo:MyID}
			SendReqChan <- statInfo
			return
		}
	}
	if _,e := simpleMsgRec[identifier]; !e {
		simpleMsgRec[identifier] = true
		for i := 0; i < total; i++ {
			m := MSGStruct{Header:MSG, Id:message.GetId(), SenderId:MyID, Data: message.GetData(), Round: message.GetRound()}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}
		//reliable accept
		//fmt.Println("Hey")
		stats := StatStruct{Id:message.GetId(), Round: message.GetRound(), Header:Stat}
		statInfo :=PrepareSend{M:stats, SendTo:MyID}
		SendReqChan <- statInfo
	}
}