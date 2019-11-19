package HRBAlgorithm

import (
	"fmt"
	"strconv"
)

/*
Note than Broadcast, MessageHandler and EchoHandler are the same for Complex version
 */
func ComplexECBroadCast(s string) {
	SimpleECBroadCast(s)
}

func ComplexECMessageHandler(m Message) {
	data,_ := m.(MSGStruct)

	identifier := data.GetId() + strconv.Itoa(data.GetRound())

	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true

		hashStr := data.GetHashData()

		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		if  c ,ok := ecDataSet[hashStr]; !ok {
			codes = make([][]byte, total)
		} else {
			codes = c
		}

		/*
			Get my corresponding code since MSG is always my corresponding code
		*/
		//Convert it to Bytes;
		codes[serverMap[MyID]], _= ConvertStringToBytes(data.GetData())

		ecDataSet[hashStr] = codes
		fmt.Println("EcDataSet, hashStr ", ecDataSet[hashStr], hashStr)

		//Main logic
		//Note data.GetData() is the string version of my code
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:hashStr, Data: data.GetData(), Round: data.GetRound(), SenderId:MyID}

		//Send Echo to all servers
		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
			sendReq := PrepareSend{M: m, SendTo:"all"}
			SendReqChan <- sendReq
		}

		checkM := MSGStruct{Id:data.GetId(), HashData:hashStr, Round:data.GetRound()}
		complexECCheck(checkM)
	}
}

func ComplexECEchoHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		echo := ECHOStruct{Header:ECHO, Id:m.GetId(), HashData:m.GetHashData(), Round: m.GetRound()}
		//fmt.Printf("%+v\n",m)
		if count, ok := EchoRecCountSet[echo]; ok {
			EchoRecCountSet[echo] = count + 1
		} else {
			EchoRecCountSet[echo] = 1
		}

		var codes [][]byte
		//include the data with key the original data and val its hashvalue
		fmt.Println("Get hashData ", echo.GetHashData())
		if  c ,ok := ecDataSet[m.GetHashData()]; !ok {
			codes = make([][]byte, total)
		} else {
			fmt.Println("It exists")
			codes = c
			fmt.Println("The codes look like", codes)
		}
		//Echo Message is always the codes of the sender

		fmt.Println(m.GetSenderId(), serverMap[m.GetSenderId()])
		codes[serverMap[m.GetSenderId()]], _ = ConvertStringToBytes(m.GetData())


		dataExist, _:= checkDataExist(echo.GetHashData())
		if !dataExist && EchoRecCountSet[echo] >= faulty + 1 {
			vals := permutation(codes, faulty + 1, total - (faulty + 1))
			fmt.Println("Receive more than faulty + 1 and the list of permutations is ", vals)
			for _, v := range vals {
				expectedHash,_ := ConvertStringToBytes(echo.GetHashData())
				inputHash, _:= ConvertStringToBytes(v)
				if compareHash(expectedHash, inputHash) {
					fmt.Println("Include ", v)
					DataSet[v] = echo.GetHashData()
					break
				}
			}
		}

		checkM := MSGStruct{Id:m.GetId(), HashData:m.GetHashData(), Round:m.GetRound()}
		complexECCheck(checkM)
	}
}

func ComplexECReqHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()

	if _, seen := ReqReceiveSet[identifier]; !seen {
		ReqReceiveSet[identifier] = true
		if exist, data := checkDataExist(m.GetHashData()); exist {
			//Decode the data
			fmt.Println("Inside req ", data)
			shards := Encode(data, faulty + 1, total - (faulty + 1))

			code := ConvertBytesToString(shards[serverMap[MyID]])
			//send val to the requested id
			fwdSendMsg := FWDStruct{Header:FWD, Id: m.GetId(), Round:m.GetRound(), SenderId:MyID, Data: code}
			req := PrepareSend{M: fwdSendMsg, SendTo: m.GetSenderId()}
			SendReqChan <- req
		}
	}

}

func complexECFwdHandler(m Message) {
	identifier := m.GetId() + strconv.Itoa(m.GetRound()) + m.GetSenderId();

	fwdIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	//fmt.Printf("Fwd: %+v\n",data)

	hashSentIdentifier := m.GetId() + strconv.Itoa(m.GetRound())
	hashStr := reqSentHash[hashSentIdentifier]
	req := REQStruct{Header:REQ, Id:m.GetId(), Round: m.GetRound(), SenderId:MyID, HashData:hashStr}

	if hasSent(ReqSentSet[req], m.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			FwdReceiveSet[identifier] = true
			//Include the code

			if count, exist := FwdRecCountSet[fwdIdentifier]; !exist {
				FwdRecCountSet[fwdIdentifier] = 1
			} else {
				FwdRecCountSet[fwdIdentifier] = count + 1
			}

			if _, e := reqSentHash[hashStr]; e {
				checkM := MSGStruct{Id:m.GetId(), HashData:hashStr, Round:m.GetRound()}
				complexECCheck(checkM)
			}

			var codes [][]byte
			//include the data with key the original data and val its hashvalue
			if  c ,ok := ecDataSet[hashStr]; !ok {
				codes = make([][]byte, total)
			} else {
				codes = c
			}

			/*
				Get my corresponding code since MSG is always my corresponding code
			*/
			//Convert it to Bytes;
			codes[serverMap[m.GetSenderId()]], _= ConvertStringToBytes(m.GetData())

			ecDataSet[hashStr] = codes
			fmt.Println("EcDataSet, hashStr, size of Foward ", ecDataSet[hashStr], hashStr, FwdRecCountSet[fwdIdentifier])

			if FwdRecCountSet[fwdIdentifier] >= faulty + 1 {
				vals := permutation(codes, faulty + 1, total - (faulty + 1))
				fmt.Println("Receive more than faulty + 1 and the list of permutations is ", vals)
				for _, v := range vals {
					expectedHash,_ := ConvertStringToBytes(hashStr)
					inputHash, _:= ConvertStringToBytes(v)
					if compareHash(expectedHash, inputHash) {
						fmt.Println("Include ", v)
						DataSet[v] = hashStr
						break
					}
				}
			}

			//check
			checkM := MSGStruct{Id:m.GetId(), HashData:hashStr, Round:m.GetRound()}
			complexECCheck(checkM)
		}
	}
}

func complexECAccHandler(m Message) {
	identifier := m.GetId()+ strconv.Itoa(m.GetRound()) + m.GetSenderId()

	if _, seen := AccReceiveSet[identifier]; !seen {
		AccReceiveSet[identifier] = true

		accM := ACCStruct{Header: ACC, Id: m.GetId(), Round: m.GetRound()}

		if l, ok := AccRecCountSet[accM]; !ok {
			l = []string{m.GetSenderId()}
			AccRecCountSet[accM] = l
		} else {
			l = append(l, m.GetSenderId())
			AccRecCountSet[accM] = l
		}

		reqMes := REQStruct{Header: REQ, Id: m.GetId(), HashData: m.GetHashData(), Round: m.GetRound(), SenderId: MyID}

		//mt.Printf("Acc in AccHandler %+v, %+v \n", accM, AccRecCountSet[accM])
		if len(AccRecCountSet[accM]) == faulty + 1 {
			//Save the Hash value that you sent request for
			sentHashIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())
			reqSentHash[sentHashIdentifier] = m.GetHashData()
			if exist, _ := checkDataExist(m.GetHashData()); !exist {
				ReqSentSet[reqMes] = AccRecCountSet[accM]
				//Send Req to these f + 1
				for _, id := range ReqSentSet[reqMes] {
					reqMessage := PrepareSend{M: reqMes, SendTo: id}
					SendReqChan <- reqMessage
				}
			}

		} else if len(AccRecCountSet[accM]) > faulty + 1 {
			//send the request to the current id
			if exist, _ := checkDataExist(m.GetHashData()); !exist {
				//Todo: Need to set a flag that if we have accepted this round, we donot need to send the req again
				sentHashIdentifier := m.GetId()+ strconv.Itoa(m.GetRound())
				exist, _ = checkDataExist(reqSentHash[sentHashIdentifier])
				if ! exist {
					l := ReqSentSet[reqMes]
					l = append(l, m.GetSenderId())
					ReqSentSet[reqMes] = l
					//send request
					reqMessage := PrepareSend{M: reqMes, SendTo: m.GetSenderId()}
					SendReqChan <- reqMessage
				}
			}
		}
		checkM := MSGStruct{Id:m.GetId(), HashData:m.GetHashData(), Round:m.GetRound()}
		complexECCheck(checkM)
	}
}

func complexECCheck(m Message) {
	fmt.Println("Inside check")

	if exist, value := checkDataExist(m.GetHashData()); exist {
		echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
		acc := ACCStruct{Header:ACC, Round:m.GetRound(), Id:m.GetId()}


		identifier := m.GetId() + strconv.Itoa(m.GetRound())

		if EchoRecCountSet[echo] >= faulty + 1 {
			if _, sent := EchoSentSet[identifier]; !sent {
				fmt.Println("Receive more than faulty + 1 echo message")
				EchoSentSet[identifier] = true
				fmt.Println("Sent Echo to other servers")
				shards := Encode(value, faulty + 1, total - (faulty + 1))
				code := ConvertBytesToString(shards[serverMap[MyID]])
				echoSend := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID, Data:code}
				sendReq := PrepareSend{M: echoSend, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		if EchoRecCountSet[echo] >= total - faulty {
			if _, sent := AccSentSet[identifier]; !sent {
				fmt.Println("Receive more than total - faulty echo message")
				AccSentSet[identifier] = true

				accSend := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				//fmt.Println("Acc sender ID "+ MyID + "," + acc.SenderId)
				//fmt.Printf("Sent Acc to other servers %+v \n" , acc)
				sendReq := PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		fmt.Println("Check Acc ", len(AccRecCountSet[acc]), AccRecCountSet[acc])

		if len(AccRecCountSet[acc]) >= faulty + 1 {
			if _,sent := AccSentSet[identifier]; !sent {
				AccSentSet[identifier] = true
				accSend := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				sendReq := PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		if len(AccRecCountSet[acc]) >= total - faulty {
			//accept
			fmt.Println("Reliable Accept " + value)
			acceptData[value] = true
		}
	}
}
