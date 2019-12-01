package HRBAlgorithm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

type digestStruct struct{
	SenderId string
	DigestM string
	Key string
}

func broadcast(s string, round int) {
	for i := 0; i < total; i++ {
		m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: s, Round:round}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func encrypt(data,key string) string{
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}


func receiveBroadCast(m Message) {
	//The sender node does not need to send Hash
	if MyID != m.GetSenderId() {
		identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound());
		digestSourceData[identifier] = m.GetData()

		data := m.GetData()
		round := m.GetRound()

		senderIndex := serverMap[MyID]
		for i := 0; i < total; i++ {
			digestM := encrypt(data, genKey)
			m := ECHOStruct{Header:ECHO, Id:m.GetId(), SenderId:MyID, HashData: genKey, Data: digestM, Round:round}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
			receiveIndex := serverMap[serverList[i]]

			dataStruct := digestStruct{Key: genKey, SenderId:MyID, DigestM:digestM}
			if arr , ok := digestRecSend[identifier]; !ok {
				arrays := make([][]digestStruct, total)
				for i := range arrays {
					arrays[i] = make([]digestStruct, total)
				}
				arrays[senderIndex][receiveIndex] = dataStruct;
				digestRecSend[identifier] = arrays
			} else {
				arr[senderIndex][receiveIndex] = dataStruct;
				digestRecSend[identifier] = arr;
			}

		}
	}
}

func receiveDigest(m Message) {
	identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound())
	digestData := digestStruct{SenderId:m.GetSenderId(), DigestM:m.GetData(), Key:m.GetHashData()}

	senderIndex := serverMap[m.GetSenderId()]
	receiverIndex := serverMap[MyID]

	var arrays [][]digestStruct
	if l,ok := digestDataMap[identifier]; !ok {
		l := []digestStruct{digestData}
		digestDataMap[identifier] = l

		if _, ok := digestRecSend[identifier]; !ok {
			arrays = make([][]digestStruct, total)
			for i := range arrays {
				arrays[i] = make([]digestStruct, total)
			}
		} else {
			arrays = digestRecSend[identifier]
		}

	} else {
		l = append(l, digestData)
		digestDataMap[identifier] = l;
		arrays = digestRecSend[identifier]
		if len(l) == total - 1 {
			data := digestSourceData[identifier]
			if validate(data, l) {
				fmt.Println("Accept Data", data)
				acceptData[data] = true
				broadcastBinary(false)
			} else {
				//Broadcast Binary
				broadcastBinary(true)
			}
		}
	}

	arrays[senderIndex][receiverIndex] = digestData
	digestRecSend[identifier] = arrays
}

func validate(targetMessage string, l []digestStruct) bool{
	for _, digestData := range l {
		key := digestData.Key
		digestM := digestData.DigestM
		targetM := encrypt(targetMessage, key)
		if targetM != digestM {
			return false
		}
	}
	return true
}

func broadCastSendRec() {

}

func recBinary() {

}

func findFaulty() {

}

func broadcastBinary(detect bool) {

}