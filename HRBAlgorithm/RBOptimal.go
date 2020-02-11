package HRBAlgorithm

import (
	"fmt"
	"strconv"
	"time"
)

var codeSet map[string] [][]byte //ecDataSet
var msgSet map[string] []string
var hashSet map[string] []string

func InitOptimal() {
	msgSet = make(map[string] []string)
	hashSet = make(map[string] []string)
}

func OptimalBroadcast(length, round int) {
	time.Sleep(3*time.Second)
	for r := 0; r < round; r ++ {
		//Generate the message
		data := RandStringBytes(length)

		/*
		Broadcast hash to other servers
		 */
		dataBytes, _ := ConvertStringToBytes(data)
		hashData := Hash(dataBytes)

		hashTag := Binary{Header:BIN, Round:r, Id:MyID, HashData:ConvertBytesToString(hashData), SenderId:MyID}
		//Broadcast HashMessage to other Servers
		for i := 0; i < total; i++ {
			sendReq := PrepareSend{M: hashTag, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}

		/*
		Broadcast codes to other Server
		 */
		var shards [][]byte
		//Generate the Coded Message
		if faulty == 0 {
			shards = Encode(data, total / 2, total - (total / 2))
		} else {
			shards = Encode(data, total - 3*faulty, total - (total - 3*faulty))
		}

		for i := 0; i < total; i++ {
			code := ConvertBytesToString(shards[i])
			m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: code, Round:r}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}
	}
}

func optimalMsgHandler(m Message) {
	data,_ := m.(MSGStruct)

	identifier := data.GetId() + strconv.Itoa(data.GetRound())

	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true

		stats := Stats{}
		stats.Start = time.Now()
		statsRecord[identifier] = stats
		fmt.Printf("Begin Stats: %+v\n",stats)


		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		if  c ,ok := ecDataSet[identifier]; !ok {
			codes = make([][]byte, total)
		} else {
			codes = c
		}

		/*
			Get my corresponding code since MSG is always my corresponding code
		*/
		//Convert it to Bytes;
		codes[serverMap[MyID]], _= ConvertStringToBytes(data.GetData())

		ecDataSet[identifier] = codes


		//Main logic
		//Note data.GetData() is the string version of my code
		//DATA IS code
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), Data: data.GetData(), Round: data.GetRound(), SenderId:MyID}

		//Send Echo to all servers
		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
			for i := 0; i < total; i++ {
				sendReq := PrepareSend{M: m, SendTo:serverList[i]}
				SendReqChan <- sendReq
			}
		}
	}
}

func optimalEchoHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()

	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		echo := ECHOStruct{Header:ECHO, Id:m.GetId(), Round: m.GetRound()}
		//fmt.Printf("%+v\n",m)
		if count, ok := EchoRecCountSet[echo]; ok {
			EchoRecCountSet[echo] = count + 1
		} else {
			EchoRecCountSet[echo] = 1
		}

		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		codeIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

		if  c ,ok := ecDataSet[codeIdentifier]; !ok {
			codes = make([][]byte, total)
		} else {
			codes = c
		}
		//Echo Message is always the codes of the sender
		codes[serverMap[m.GetSenderId()]], _ = ConvertStringToBytes(m.GetData())
		ecDataSet[codeIdentifier] = codes

		if EchoRecCountSet[echo] >= total - faulty {
			var vals []string
			if faulty == 0 {
				vals = permutation(codes, total / 2, total - (total / 2))
			} else {
				vals = permutation(codes, total - 3 * faulty, total - (total - 3 * faulty))
			}

			fmt.Println("Decode the messageSet: " , vals)
			if checkDecodeSet(vals) {
				if _, e := msgSet[codeIdentifier]; !e {
					msgSet[codeIdentifier] = []string{vals[0]}
				} else {
					if !checkValExist(vals[0], msgSet[codeIdentifier]) {
						msgSet[codeIdentifier] = append(msgSet[codeIdentifier], vals[0])
					}
				}

				exist, hashData := checkHashExist(vals[0] ,codeIdentifier)
				fmt.Println(exist, hashData)
				if exist {
					accStruct := ACCStruct{Header:ACC, HashData:hashData, Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

					if _, e := AccSentSet[codeIdentifier]; !e {
						AccSentSet[codeIdentifier] = true
						for i := 0; i < total; i++ {
							sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
							SendReqChan <- sendReq
						}
					}

				}
			}
		}
	}
}

func optimalAccHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()

	roundIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

	if _, seen := AccReceiveSet[identifier]; !seen {
		AccReceiveSet[identifier] = true

		accM := ACCStruct{Header: ACC, Id: m.GetId(), Round: m.GetRound(), HashData:m.GetHashData()}

		if l, ok := AccRecCountSet[accM]; !ok {
			l = []string{m.GetSenderId()}
			AccRecCountSet[accM] = l
		} else {
			l = append(l, m.GetSenderId())
			AccRecCountSet[accM] = l
		}

		if len(AccRecCountSet[accM]) == faulty + 1 {
			if _, e:=AccSentSet[roundIdentifier]; !e {
				AccSentSet[roundIdentifier] = true
				accM = ACCStruct{Header: ACC, Id: m.GetId(), Round: m.GetRound(), HashData:m.GetHashData(), SenderId:MyID}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M: accM, SendTo:serverList[i]}
					SendReqChan <- sendReq
				}
			}
		}
		optimalCheck(m)
	}
}

func optimalReqHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()
	roundIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

	if _, seen := ReqReceiveSet[identifier]; !seen {
		ReqReceiveSet[identifier] = true
		exist, val := checkOptimalDataExist(m.GetHashData(), roundIdentifier)
		if exist {
			fwd := FWDStruct{Header:FWD, Id:m.GetId(), SenderId:m.GetSenderId(), Round:m.GetRound(), Data:val, HashData:m.GetHashData()}
			sendReq := PrepareSend{M: fwd, SendTo: m.GetSenderId()}
			SendReqChan <- sendReq
		}
	}
}

func optimalFwdHandler(m Message) {
	identifier := m.GetId() + strconv.Itoa(m.GetRound()) + m.GetSenderId();

	roundIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	//fmt.Printf("Fwd: %+v\n",data)


	hashStr := reqSentHash[roundIdentifier]
	req := REQStruct{Header:REQ, Id:m.GetId(), Round: m.GetRound(), SenderId:MyID, HashData:hashStr}

	if hasSent(ReqSentSet[req], m.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			FwdReceiveSet[identifier] = true
			dataBytes, _ := ConvertStringToBytes(m.GetData())
			if ConvertBytesToString(Hash(dataBytes)) == m.GetHashData() {

				//Send Acc
				exist, hashData := checkHashExist(m.GetData() ,roundIdentifier)
				if exist {
					accStruct := ACCStruct{Header:ACC, HashData:hashData, Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

					if _, e := AccSentSet[roundIdentifier]; !e {
						AccSentSet[roundIdentifier] = true
						for i := 0; i < total; i++ {
							sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
							SendReqChan <- sendReq
						}
					}

				}

				if vals, e := msgSet[roundIdentifier]; !e {
					msgSet[roundIdentifier] = []string{m.GetData()}
				} else {
					msgSet[roundIdentifier] = append(vals, m.GetData())
				}
				optimalCheck(m)
			}
		}
	}
}

func optimalHashTagHandler(m Message) {
	fmt.Printf("Fwd: %+v\n",m)
	roundIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

	if _, e := hashSet[roundIdentifier]; e {
		hashSet[roundIdentifier] = append(hashSet[roundIdentifier], m.GetHashData())
	} else {
		hashSet[roundIdentifier] = []string{m.GetHashData()}
	}

	fmt.Println("Hashset:" + hashSet[roundIdentifier][0] + " " + roundIdentifier)

	//Send Acc
	exist, _ := checkOptimalDataExist(m.GetHashData(), roundIdentifier)

	if exist {
		accStruct := ACCStruct{Header:ACC, HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

		if _, e := AccSentSet[roundIdentifier]; !e {
			AccSentSet[roundIdentifier] = true
			for i := 0; i < total; i++ {
				sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
				SendReqChan <- sendReq
			}
		}

	}
}

func optimalCheck(m Message) {
	roundIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	seenMap := make(map[string] bool)
	for _, x := range hashSet[roundIdentifier] {
		if _, e := seenMap[x]; !e {
			seenMap[x] = true
			accM := ACCStruct{Header: ACC, Id: m.GetId(), Round: m.GetRound(), HashData:x}
			if len(AccRecCountSet[accM]) >= total - faulty {
				exist, data := checkOptimalDataExist(x, roundIdentifier)
				if exist {
					fmt.Println("Reliable Accept", data)
				} else {
					req := REQStruct{Header:REQ, Id:m.GetId(), Round: m.GetRound(), SenderId:MyID, HashData:x}
					ReqSentSet[req] = AccRecCountSet[accM]
					for _, server := range AccRecCountSet[accM] {
						sendReq := PrepareSend{M: req, SendTo: server}
						SendReqChan <- sendReq
					}
				}
			}

		}
	}
}

func checkOptimalDataExist(hashData, roundIdentifier string) (bool, string){
	vals := msgSet[roundIdentifier]
	for i := 0; i < len(vals); i++ {
		bytes, _ := ConvertStringToBytes(vals[i])
		dataBytes := Hash(bytes)
		dataHashStr := ConvertBytesToString(dataBytes)
		if dataHashStr == hashData {
			return true, vals[i]
		}
	}
	return false, ""
}

//input is data
func checkHashExist(m, identifier string) (bool, string){
	messageBytes, _ := ConvertStringToBytes(m)
	hashData := Hash(messageBytes)
	hashDataString := ConvertBytesToString(hashData)
	fmt.Println(m + ", " + hashDataString + ", " + identifier)

	hashes := hashSet[identifier]
	for i := 0; i < len(hashes); i++ {
		fmt.Println("Print checkHash exsist", hashes[i], hashDataString)
		if hashes[i] == hashDataString {
			return true, hashes[i]
		}
	}
	return false, ""
}


func checkDecodeSet(vals []string) bool{
	target := vals[0]
	for i := 0; i < len(vals); i++ {
		if vals[i] != target {
			return false
		}
	}
	return true
}

func checkValExist(input string, vals []string) bool{
	for _, val := range vals {
		if val == input {
			return true
		}
	}
	return false
}