package HRBAlgorithm

import (
	"fmt"
	"strconv"
)

var ByzCodeCounter map[string] int
var ByzCodeElement map[string] [][]byte

func InitByzCode() {
	ByzCodeCounter = make (map[string] int)
	ByzCodeElement = make (map[string] [][]byte)
}

func ECByzBroadCast(s string, round int) {
	//need to make sure that coded element > f
	var shards[][] byte
	if faulty == 0 {
		shards = Encode(s, total, total)
	} else {
		shards = Encode(s, total - faulty, 2*(total) - (total - faulty))
	}
	fmt.Println("Shards are ", shards)


	for i := 0; i < total; i++ {
		code1 := ConvertBytesToString(shards[i])
		code2 := ConvertBytesToString(shards[i + total - 1])
		m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: code1, Round: round, HashData: code2}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func ByzRecMsg(m Message) {
	identifier := identifierCreate(m.GetId(), m.GetRound())
	count, exist := ByzCodeCounter[identifier]

	if exist {
		ByzCodeCounter[identifier] = count + 2
	} else {
		ByzCodeCounter[identifier] = 2
		ByzCodeElement[identifier] = make([][]byte, 2 * total)
	}

	index1 := serverMap[MyID]
	index2 := index1 + total - 1
	code1, _ := ConvertStringToBytes(m.GetData())
	code2, _ := ConvertStringToBytes(m.GetHashData())
	fmt.Println("code1: ", code1, " code2: ", code2)
	ByzCodeElement[identifier][index1] = code1
	ByzCodeElement[identifier][index2] = code2

	code := ConvertBytesToString(code1)
	id := m.GetId();
	round := m.GetRound()
	//Send Echo
	for i := 0; i < total; i++ {
		message := ECHOStruct{Header:ECHO, Id:id, SenderId:MyID, Data: code, Round: round}
		sendReq := PrepareSend{M: message, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func ByzRecEcho(m Message) {
	if MyID != m.GetSenderId() {
		identifier := identifierCreate(m.GetId(), m.GetRound())
		count, exist := ByzCodeCounter[identifier]

		if exist {
			ByzCodeCounter[identifier] = count + 1
		} else {
			ByzCodeCounter[identifier] = 1
			ByzCodeElement[identifier] = make([][]byte, 2 * total)
		}

		index := serverMap[m.GetSenderId()]
		code, _ := ConvertStringToBytes(m.GetData())
		ByzCodeElement[identifier][index] = code

		if ByzCodeCounter[identifier] == total {
			var vals []string
			if faulty == 0 {
				vals = permutation(ByzCodeElement[identifier], total, total)
			} else {
				vals = permutation(ByzCodeElement[identifier], total - faulty, 2*(total) - (total - faulty))
			}
			fmt.Println(ByzCodeCounter[identifier], ByzCodeElement[identifier])
			detected := validateByzCode(vals)
			broadcastBinary(detected, m.GetId(), m.GetRound())
		}
	}
}

func ByzRecBin(m Message) {
	if m.GetSenderId() != MyID {
		identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound())
		if l, ok := binarySet[identifier]; !ok {
			firstL := []Message{m}
			binarySet[identifier] = firstL
		} else {
			l = append(l, m)
			binarySet[identifier] = l
			fmt.Println(binarySet[identifier])
			if len(l) == total - 1 {
				detect := checkDetect(l)
				if detect {
					fmt.Println("Fail to Accept Data")
				} else {
					data := digestSourceData[identifier]
					fmt.Println("Accept Data", data)
				}
			}
		}
	}
}

func validateByzCode(vals []string) bool{
	fmt.Println(vals)
	data := vals[0]
	for _, val := range vals {
		if val != data {
			fmt.Println("Wrong", val, data)
			return true
		}
	}
	return false
}



