package HRBAlgorithm

import (
	"strconv"
	"time"
)

type brachaCountStruct struct{
	Id string
	Round int
	Data string
}

var (
	brachaMsgRecSet map[string] bool
	brachaEchoRecSet map[string] bool
	brachaAccRecSet map[string] bool
	brachaDataRecSet map[string] string

	brachaEchoHasSent map[string] bool
	brachaAccHasSent map[string] bool

	brachaEchoCount map[brachaCountStruct] int
	brachaAccCount map[brachaCountStruct] int

	brachaAccept map[string] bool
)

func Initbracha(round int) {
	brachaMsgRecSet = make (map[string] bool, round / 2)
	brachaEchoRecSet = make (map[string] bool, round / 2)
	brachaAccRecSet = make (map[string] bool, round / 2)
	brachaDataRecSet = make (map[string] string, round / 2)

	brachaEchoHasSent = make(map[string] bool, round / 2)
	brachaAccHasSent = make(map[string] bool, round / 2)

	brachaEchoCount = make(map[brachaCountStruct] int, round / 2)
	brachaAccCount = make (map[brachaCountStruct] int, round / 2)

	brachaAccept = make (map[string] bool, round / 2)
}

func BrachaBroadCast(length , round int) {
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

func senderIdentifierCreate(id string, round int, senderId string) string{
	return senderId + ":" + id + ":" + strconv.Itoa(round)
}

func roundIdentifierCreate(id string, round int) string{
	return id + ":" + strconv.Itoa(round)
}

func brachaMessageHandler(message Message) {
	id := message.GetId()
	round := message.GetRound()
	senderId := message.GetSenderId()
	roundIdentifier := roundIdentifierCreate(id, round)
	data := message.GetData()

	if senderId == id {
		if _, e := brachaMsgRecSet[roundIdentifier]; !e {
			brachaMsgRecSet[roundIdentifier] = true

			if _, e := brachaEchoHasSent [roundIdentifier]; !e {
				brachaEchoHasSent [roundIdentifier] = true

				EchoStruct := ECHOStruct{Id:id, SenderId:MyID, Round: round, Data:data, Header:ECHO}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M:EchoStruct, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}
			}
		}
	}
}

func brachaEchoHandler(message Message) {
	id := message.GetId()
	round := message.GetRound()
	senderId := message.GetSenderId()
	roundIdentifier := roundIdentifierCreate(id, round)
	senderIdentifier := senderIdentifierCreate(id, round, senderId)
	data := message.GetData()

	if _, e := brachaEchoRecSet[senderIdentifier]; !e {
		brachaEchoRecSet[senderIdentifier] = true
		EchoCountStruct := brachaCountStruct{Id:id, Round:round, Data:data}

		if count, e := brachaEchoCount[EchoCountStruct]; !e {
			brachaEchoCount[EchoCountStruct] = 1
		} else {
			brachaEchoCount[EchoCountStruct] = count + 1
		}

		if brachaEchoCount[EchoCountStruct] == (total + faulty) / 2 {

			if _, hasSent := brachaEchoHasSent[roundIdentifier]; !hasSent {
				brachaEchoHasSent[roundIdentifier] = true
				EchoStruct := ECHOStruct{Id:id, SenderId:MyID, Round: round, Data:data, Header:ECHO}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M:EchoStruct, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}
			}

			if _, hasSent := brachaAccHasSent[roundIdentifier]; !hasSent {
				brachaAccHasSent[roundIdentifier] = true
				ACCStruct := ACCStruct{Id:id, SenderId:MyID, Round:round, Data: data, Header:ACC}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M:ACCStruct, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}

			}
		}
	}
}

func brachaAccHandler(message Message) {
	id := message.GetId()
	round := message.GetRound()
	senderId := message.GetSenderId()
	roundIdentifier := roundIdentifierCreate(id, round)
	senderIdentifier := senderIdentifierCreate(id, round, senderId)
	data := message.GetData()

	if _, e := brachaAccRecSet[senderIdentifier ]; !e {
		brachaAccRecSet[senderIdentifier ] = true

		accCountStruct := brachaCountStruct{Id:id, Round:round, Data:data}

		if count, e := brachaAccCount[accCountStruct]; !e {
			brachaAccCount[accCountStruct] = 1
		} else {
			brachaAccCount[accCountStruct] = count + 1
		}

		if brachaAccCount[accCountStruct] == faulty + 1 {

			if _, hasSent := brachaAccHasSent[roundIdentifier]; !hasSent {
				brachaAccHasSent[roundIdentifier] = true

				EchoStruct := ECHOStruct{Id:id, SenderId:MyID, Round: round, Data:data, Header:ECHO}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M:EchoStruct, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}
			}

			if _, hasSent := brachaAccHasSent[roundIdentifier]; !hasSent {
				brachaAccHasSent[roundIdentifier] = true

				ACCStruct := ACCStruct{Id:id, SenderId:MyID, Round:round, Data: data, Header:ACC}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M:ACCStruct, SendTo: serverList[i]}
					SendReqChan <- sendReq
				}
			}
		}

		checkBrachaReliableAccept(accCountStruct)
	}
}

func checkBrachaReliableAccept(countStruct brachaCountStruct) {
	roundIdentifier := roundIdentifierCreate(countStruct.Id, countStruct.Round)
	if brachaAccCount[countStruct] == 2 * faulty + 1 {
		if _, e := brachaAccept[roundIdentifier]; !e {
			brachaAccept[roundIdentifier] = true
			//fmt.Println("Relable accept", countStruct.Data)
			stats := StatStruct{Id:countStruct.Id, Round: countStruct.Round, Header:Stat}
			statInfo :=PrepareSend{M:stats, SendTo:MyID}
			SendReqChan <- statInfo
		}
	}

}
