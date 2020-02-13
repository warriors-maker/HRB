package HRBAlgorithm

import (
	"strconv"
	"time"
)

var codeSet map[string] [][]byte //ecDataSet
var msgSet map[string] []string
var hashOptimalSet map[string] []string
var messageOptimalReceiveSet map[string] bool
var echoOptimalSentSet map[string] bool
var echoOptimalReceiveSet map[string] bool
var echoOptimalRecCountSet map[ECHOStruct] int
var accOptimalSentSet map[string] bool
var accOptimalRecCountSet map[ACCStruct] []string
var reqOptimalReceiveSet map[string] bool
var reqOptimalSentHash map[string] string
var fwdOptimalReceiveSet map[string] bool
var reqOptimalSentSet map[REQStruct] []string
var OptimalSendReqChan chan PrepareSend
var optimalReliableAccept map[string] bool

func InitOptimal(round int) {
	msgSet = make(map[string] []string, round / 2)
	hashOptimalSet = make(map[string] []string, round / 2)
	messageOptimalReceiveSet = make(map[string] bool, round / 2)
	codeSet = make (map[string] [][]byte, round / 2)
	echoOptimalSentSet = make(map[string] bool, round / 2)
	echoOptimalReceiveSet = make(map[string] bool, round / 2)
	echoOptimalRecCountSet = make(map[ECHOStruct] int, round / 2)
	accOptimalSentSet = make (map[string] bool, round / 2)
	accOptimalRecCountSet = make(map[ACCStruct] []string, round/2)
	reqOptimalReceiveSet = make(map[string] bool, round / 2)
	reqOptimalSentHash = make(map[string] string, round / 2)
	fwdOptimalReceiveSet = make(map[string] bool, round / 2)
	reqOptimalSentSet = make(map[REQStruct] []string, round / 2)
	optimalReliableAccept = make(map[string] bool, round / 2)
	OptimalSendReqChan = make(chan PrepareSend)
}

func hashTagBroadcast(round int, data string) {
	//if i % 200 == 0 {
	//	time.Sleep(1*time.Second)
	//}
	m := MSGStruct{Id: MyID, SenderId:MyID, Data: data, Header:MSG, Round:round}
	for _, server := range serverList {
		//fmt.Println("Protocal send to ", server)
		sendReq := PrepareSend{M: m, SendTo: server}
		SendReqChan <- sendReq
	}
}

func OptimalBroadcast(length, round int) {
	time.Sleep(3*time.Second)
	for r := 0; r < round; r ++ {
		//Generate the message
		data := RandStringBytes(length)

		/*
			Using Chocolate Broadcast
		*/


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
		//Chocolate Broadcast
		hashTagBroadcast(r, data)

		for i := 0; i < total; i++ {
			// Chocolate Broadcast
			//Coded Broadcast
			code := ConvertBytesToString(shards[i])
			m := MSGStruct{Header:MSG_OPT, Id:MyID, SenderId:MyID, Data: code, Round:r}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			OptimalSendReqChan <- sendReq
		}
	}
}

func optimalMsgHandler(m Message) {
	data,_ := m.(MSGStruct)

	identifier := data.GetId() + strconv.Itoa(data.GetRound())

	if _, seen := messageOptimalReceiveSet[identifier]; !seen {
		messageOptimalReceiveSet[identifier] = true


		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		if  c ,ok := codeSet[identifier]; !ok {
			codes = make([][]byte, total)
		} else {
			codes = c
		}

		/*
			Get my corresponding code since MSG is always my corresponding code
		*/
		//Convert it to Bytes;
		codes[serverMap[MyID]], _= ConvertStringToBytes(data.GetData())

		codeSet[identifier] = codes


		//Main logic
		//Note data.GetData() is the string version of my code
		//DATA IS code
		m := ECHOStruct{Header:ECHO_OPT, Id:data.GetId(), Data: data.GetData(), Round: data.GetRound(), SenderId:MyID}

		//Send Echo to all servers
		if _, sent := echoOptimalSentSet[identifier]; !sent {
			echoOptimalSentSet[identifier] = true
			for i := 0; i < total; i++ {
				sendReq := PrepareSend{M: m, SendTo:serverList[i]}
				OptimalSendReqChan <- sendReq
			}
		}
	}
}

func optimalEchoHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()

	if _,seen := echoOptimalReceiveSet[identifier]; !seen {
		echoOptimalReceiveSet[identifier] = true

		echo := ECHOStruct{Header:ECHO_OPT, Id:m.GetId(), Round: m.GetRound()}
		//fmt.Printf("%+v\n",m)
		if count, ok := echoOptimalRecCountSet[echo]; ok {
			echoOptimalRecCountSet[echo] = count + 1
		} else {
			echoOptimalRecCountSet[echo] = 1
		}

		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		codeIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

		if  c ,ok := codeSet[codeIdentifier]; !ok {
			codes = make([][]byte, total)
		} else {
			codes = c
		}
		//Echo Message is always the codes of the sender
		codes[serverMap[m.GetSenderId()]], _ = ConvertStringToBytes(m.GetData())
		codeSet[codeIdentifier] = codes

		if echoOptimalRecCountSet[echo] >= total - faulty {
			var vals []string
			if faulty == 0 {
				vals = permutation(codes, total / 2, total - (total / 2))
			} else {
				vals = permutation(codes, total - 3 * faulty, total - (total - 3 * faulty))
			}

			//fmt.Println("Decode the messageSet: " , vals)
			if checkDecodeSet(vals) {
				if _, e := msgSet[codeIdentifier]; !e {
					msgSet[codeIdentifier] = []string{vals[0]}
				} else {
					if !checkValExist(vals[0], msgSet[codeIdentifier]) {
						msgSet[codeIdentifier] = append(msgSet[codeIdentifier], vals[0])
					}
				}

				exist, hashData := checkHashExist(vals[0] ,codeIdentifier)
				//fmt.Println(exist, hashData)
				if exist {
					accStruct := ACCStruct{Header:ACC_OPT, HashData:hashData, Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

					if _, e := accOptimalSentSet[codeIdentifier]; !e {
						accOptimalSentSet[codeIdentifier] = true
						for i := 0; i < total; i++ {
							sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
							OptimalSendReqChan <- sendReq
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

	if _, seen := accOptimalSentSet[identifier]; !seen {
		accOptimalSentSet[identifier] = true

		accM := ACCStruct{Header: ACC_OPT, Id: m.GetId(), Round: m.GetRound(), HashData:m.GetHashData()}

		if l, ok := accOptimalRecCountSet[accM]; !ok {
			l = []string{m.GetSenderId()}
			accOptimalRecCountSet[accM] = l
		} else {
			l = append(l, m.GetSenderId())
			accOptimalRecCountSet[accM] = l
		}

		if len(accOptimalRecCountSet[accM]) == faulty + 1 {
			if _, e:=accOptimalSentSet[roundIdentifier]; !e {
				accOptimalSentSet[roundIdentifier] = true
				accM = ACCStruct{Header: ACC_OPT, Id: m.GetId(), Round: m.GetRound(), HashData:m.GetHashData(), SenderId:MyID}
				for i := 0; i < total; i++ {
					sendReq := PrepareSend{M: accM, SendTo:serverList[i]}
					OptimalSendReqChan <- sendReq
				}
			}
		}
		optimalCheck(m)
	}
}

func optimalReqHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()
	roundIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())

	if _, seen := reqOptimalReceiveSet[identifier]; !seen {
		reqOptimalReceiveSet[identifier] = true
		exist, val := checkOptimalDataExist(m.GetHashData(), roundIdentifier)
		if exist {
			fwd := FWDStruct{Header:FWD_OPT, Id:m.GetId(), SenderId:m.GetSenderId(), Round:m.GetRound(), Data:val, HashData:m.GetHashData()}
			sendReq := PrepareSend{M: fwd, SendTo: m.GetSenderId()}
			OptimalSendReqChan <- sendReq
		}
	}
}

func optimalFwdHandler(m Message) {
	identifier := m.GetId() + strconv.Itoa(m.GetRound()) + m.GetSenderId();

	roundIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	//fmt.Printf("Fwd: %+v\n",data)


	hashStr := reqOptimalSentHash[roundIdentifier]
	req := REQStruct{Header:REQ_OPT, Id:m.GetId(), Round: m.GetRound(), SenderId:MyID, HashData:hashStr}

	if hasSent(reqOptimalSentSet[req], m.GetSenderId()) {
		if _,seen := fwdOptimalReceiveSet[identifier]; !seen {
			fwdOptimalReceiveSet[identifier] = true
			dataBytes, _ := ConvertStringToBytes(m.GetData())
			if ConvertBytesToString(Hash(dataBytes)) == m.GetHashData() {

				//Send Acc
				exist, hashData := checkHashExist(m.GetData() ,roundIdentifier)
				if exist {
					accStruct := ACCStruct{Header:ACC_OPT, HashData:hashData, Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

					if _, e := accOptimalSentSet[roundIdentifier]; !e {
						accOptimalSentSet[roundIdentifier] = true
						for i := 0; i < total; i++ {
							sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
							OptimalSendReqChan <- sendReq
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
	//fmt.Printf("Receive from HRB: %+v\n",m)

	roundIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())
	//fmt.Println(roundIdentifier)

	if _, e := hashOptimalSet[roundIdentifier]; e {
		hashOptimalSet[roundIdentifier] = append(hashOptimalSet[roundIdentifier], m.GetHashData())
	} else {
		hashOptimalSet[roundIdentifier] = []string{m.GetHashData()}
	}

	//fmt.Println("Hashset:" + hashOptimalSet[roundIdentifier][0] + " " + roundIdentifier)

	//Send Acc
	exist, _ := checkOptimalDataExist(m.GetHashData(), roundIdentifier)

	//fmt.Println("HashTagHandler: ", exist)
	if exist {
		accStruct := ACCStruct{Header:ACC_OPT, HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID, Round:m.GetRound()}

		if _, e := accOptimalSentSet[roundIdentifier]; !e {
			accOptimalSentSet[roundIdentifier] = true
			for i := 0; i < total; i++ {
				//fmt.Println("Sent Acc: " + serverList[i])
				sendReq := PrepareSend{M: accStruct, SendTo:serverList[i]}
				OptimalSendReqChan <- sendReq
			}
		}

	}
}

func optimalCheck(m Message) {
	roundIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	seenMap := make(map[string] bool)
	for _, x := range hashOptimalSet[roundIdentifier] {
		if _, e := seenMap[x]; !e {
			seenMap[x] = true
			accM := ACCStruct{Header: ACC_OPT, Id: m.GetId(), Round: m.GetRound(), HashData:x}
			//fmt.Printf("ReliableAccept %+v \n" , accOptimalRecCountSet[accM])
			if len(accOptimalRecCountSet[accM]) >= total - faulty {
				exist, _ := checkOptimalDataExist(x, roundIdentifier)
				if exist {
					if _, e := optimalReliableAccept[roundIdentifier]; !e {
						optimalReliableAccept[roundIdentifier] = true
						//fmt.Println("Optimal Reliable Accept")
						stats := StatStruct{Id:m.GetId(), Round: m.GetRound(), Header:Stat}
						statInfo :=PrepareSend{M:stats, SendTo:MyID}
						SendReqChan <- statInfo
					}
				} else {
					req := REQStruct{Header:REQ_OPT, Id:m.GetId(), Round: m.GetRound(), SenderId:MyID, HashData:x}
					reqOptimalSentSet[req] = accOptimalRecCountSet[accM]
					for _, server := range accOptimalRecCountSet[accM] {
						sendReq := PrepareSend{M: req, SendTo: server}
						OptimalSendReqChan <- sendReq
					}
				}
			}
		}
	}
}

func checkOptimalDataExist(hashData, roundIdentifier string) (bool, string){
	vals := msgSet[roundIdentifier]
	for i := 0; i < len(vals); i++ {
		dataHashStr := ConvertBytesToString(Hash([]byte(vals[i])))
		//fmt.Println("Hey: " + dataHashStr, hashData)
		if dataHashStr == hashData {
			return true, vals[i]
		}
	}
	return false, ""
}

//input is data
func checkHashExist(m, identifier string) (bool, string){
	hashDataString := ConvertBytesToString(Hash([]byte(m)))

	hashes := hashOptimalSet[identifier]
	for i := 0; i < len(hashes); i++ {
		//fmt.Println("Print checkHash exsist", hashes[i], hashDataString)
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
