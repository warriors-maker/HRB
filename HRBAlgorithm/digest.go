package HRBAlgorithm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

var dataMap map[string] []digestStruct
var dataFromSrc map[string] string

type digestStruct struct{
	SenderId string
	DigestM string
	Key string
}

func InitDigest() {
	dataMap = make(map[string] []digestStruct)
	dataFromSrc = make(map[string] string)
}

func encrypt(data,key string) string{
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func validate(targetMessage string, l []digestStruct) bool{
	fmt.Println(targetMessage)

	for _, digestData := range l {
		if (digestStruct{}) != digestData {
			key := digestData.Key
			digestM := digestData.DigestM
			targetM := encrypt(targetMessage, key)
			fmt.Println(digestM, targetM)
			if targetM != digestM {
				return false
			}
		}
	}
	return true
}


func BroadcastPrepare(s string, round int) {
	time.Sleep(3*time.Second)
	identifier := MyID + ":" + strconv.Itoa(round);
	dataFromSrc[identifier] = s

	for i := 0; i < total; i++ {
		m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: s, Round:round}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func checkDetect(binarys [] Message) bool{
	count0 := 0
	count1 := 1
	for _, b := range binarys {
		if b.GetHashData() == "1" {
			count1 += 1
		} else {
			count0 += 1
		}
	}
	if count0 > count1 {
		return false
	} else {
		return true
	}
}

func broadcastBinary(detect bool, Id string, round int) {
	var detectString string
	if detect {
		fmt.Println("Detect")
		detectString = "1"
	} else {
		fmt.Println("Does not detect Detect")
		detectString = "0"
	}
	for i := 0; i < total; i++ {
		m := Binary{Header: BIN, SenderId:MyID, Id:Id, round:round, HashData:detectString}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

/*
Receive
 */

func receivePrepareFromSrc(m Message) {
	//The sender node does not need to send Hash
	identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound());
	dataFromSrc[identifier] = m.GetData()
	if MyID != m.GetSenderId() {
		data := m.GetData()
		round := m.GetRound()
		//Broadcast to all other servers
		for i := 0; i < total; i++ {
			//digestM is the encryption of the data
			digestM := encrypt(data, genKey)
			m := ECHOStruct{Header:ECHO, Id:m.GetId(), SenderId:MyID, HashData: genKey, Data: digestM, Round:round}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq //send to other servers

			//senderIndex := serverMap[MyID]
			//receiveIndex := serverMap[serverList[i]]
			//
			//dataStruct := digestStruct{Key: genKey, SenderId:MyID, DigestM:digestM}


			//if arr , ok := digestRecSend[identifier]; !ok {
			//	arrays := make([][]digestStruct, total)
			//	for i := range arrays {
			//		arrays[i] = make([]digestStruct, total)
			//	}
			//	arrays[senderIndex][receiveIndex] = dataStruct;
			//	digestRecSend[identifier] = arrays
			//} else {
			//	arr[senderIndex][receiveIndex] = dataStruct;
			//	digestRecSend[identifier] = arr;
			//}
		}
	}
}

func recBinary(m Message) {
	if m.GetSenderId() != MyID {
		identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound())
		if l, ok := binarySet[identifier]; !ok {
			firstL := []Message{m}
			binarySet[identifier] = firstL
		} else {
			l = append(l, m)
			binarySet[identifier] = l
			fmt.Println(binarySet[identifier])
			if len(l) == digestTrustCount - 1 {
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

func receiveDigestFromOthers(m Message) {
	identifier := m.GetId() + ":" + strconv.Itoa(m.GetRound())
	digestData := digestStruct{SenderId:m.GetSenderId(), DigestM:m.GetData(), Key:m.GetHashData()}

	digestList, exist := dataMap[identifier]

	if !exist {
		firstDigestList := make([]digestStruct, 0)
		firstDigestList = append(firstDigestList, digestData)
		dataMap[identifier] = firstDigestList
	} else {
		digestList := append(digestList, digestData)
		dataMap[identifier] = digestList
	}

	fmt.Println(dataMap[identifier], len(dataMap[identifier]))
	if len (dataMap[identifier]) == total - 1 {
		data := dataFromSrc[identifier]
		if validate(data, dataMap[identifier]) {
			fmt.Println("Broadcast detected No")
			broadcastBinary(false, m.GetId(), m.GetRound())
		} else {
			//Broadcast Binary
			fmt.Println("Broadcast detected Yes")
			broadcastBinary(true, m.GetId(), m.GetRound())
		}
	}

	//receiverIndex := serverMap[MyID]

	//var arrays [][]digestStruct
	//if l,ok := digestDataMap[identifier]; !ok {
	//	l := []digestStruct{digestData}
	//	digestDataMap[identifier] = l
	//
	//	if _, ok := digestRecSend[identifier]; !ok {
	//		arrays = make([][]digestStruct, total)
	//		for i := range arrays {
	//			arrays[i] = make([]digestStruct, total)
	//		}
	//	} else {
	//		arrays = digestRecSend[identifier]
	//	}
	//
	//} else {
	//	l = append(l, digestData)
	//	digestDataMap[identifier] = l;
	//	arrays = digestRecSend[identifier]
	//	if len(l) == total - 1 {
	//		data := digestSourceData[identifier]
	//		if validate(data, l) {
	//			broadcastBinary(false, m.GetId(), m.GetRound())
	//		} else {
	//			//Broadcast Binary
	//			broadcastBinary(true, m.GetId(), m.GetRound())
	//		}
	//	}
	//}
	//
	//arrays[senderIndex][receiverIndex] = digestData
	//digestRecSend[identifier] = arrays
}



/*
Do not need to find faulty anymore
*/

//func checkDuplicate(l []string, id string) bool{
//	for _, val := range l {
//		if val == id {
//			return true
//		}
//	}
//	return false
//}

//func broadCastSendRec(id string, round int, recSendArr [][]digestStruct) {
//	recSend := RecSend{Header: RSS, Id:id, SenderId:MyID, RecSend:recSendArr, round:round}
//	for i := 0; i < total; i++ {
//		sendReq := PrepareSend{M: recSend, SendTo: serverList[i]}
//		SendReqChan <- sendReq
//	}
//}

//func findFaulty(augmented map[string] [][]digestStruct) {
//	for senderId, recSend := range augmented {
//		//check current send with other receive
//		senderIdx := serverMap[senderId]
//
//		send_arr := recSend[senderIdx]
//
//		for receiveId, aux := range augmented {
//			receiveIdx := serverMap[receiveId]
//			sendDigest := send_arr[receiveIdx]
//			receiveDigest := aux[senderIdx][receiveIdx]
//			if sendDigest.DigestM != receiveDigest.DigestM || sendDigest.Key != receiveDigest.Key {
//				l1 := notTrustedListMap[receiveId]
//				if ! checkDuplicate(l1, senderId) {
//					l1 = append(l1, senderId)
//					notTrustedListMap[receiveId] = l1
//				}
//				l2 := notTrustedListMap[senderId]
//				if ! checkDuplicate(l2, receiveId) {
//					l2 = append(l2, receiveId)
//					notTrustedListMap[senderId] = l2
//				}
//				if senderId == MyID {
//					faultySet[receiveId] = true
//				}
//			}
//		}
//
//		//check current receive with other send
//		for i := 0; i < total; i++ {
//			receiveDigest := recSend[i][senderIdx]
//			sendDigest := augmented[serverList[i]][i][senderIdx]
//			receiveId := serverList[i]
//
//			if receiveDigest.DigestM != sendDigest.DigestM || sendDigest.Key != receiveDigest.Key {
//				l1 := notTrustedListMap[receiveId]
//				if ! checkDuplicate(l1, senderId) {
//					l1 = append(l1, senderId)
//					notTrustedListMap[receiveId] = l1
//				}
//				l2 := notTrustedListMap[senderId]
//				if ! checkDuplicate(l2, receiveId) {
//					l2 = append(l2, receiveId)
//					notTrustedListMap[senderId] = l2
//				}
//
//				if senderId == MyID {
//					faultySet[receiveId] = true
//				}
//			}
//		}
//	}
//
//	//Check whether a Server is trusted by enough people
//	for server, conflictList := range notTrustedListMap{
//		if len(conflictList) >= faulty + 1 {
//			faultySet[server] = true
//		}
//	}
//}

//func receiveRecSend(message Message) {
//	identifier := message.GetId() + ":" + strconv.Itoa(message.GetRound())
//	recSendStruct, ok := message.(RecSend)
//	if ok {
//		if existMap, exist := augmentRecSend[identifier]; !exist {
//			m := make(map[string] [][]digestStruct)
//			m[message.GetSenderId()] = recSendStruct.RecSend
//			augmentRecSend[identifier] = m
//		} else {
//			existMap[message.GetSenderId()] = recSendStruct.RecSend
//			augmentRecSend[identifier] = existMap
//		}
//	}
//
//	if len(augmentRecSend[identifier]) == digestTrustCount {
//		findFaulty(augmentRecSend[identifier])
//	}
//}