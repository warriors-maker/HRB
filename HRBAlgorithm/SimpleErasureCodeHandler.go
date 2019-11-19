package HRBAlgorithm

import (
	"fmt"
	"strconv"
)

func SimpleECBroadCast(s string) {
	shards := Encode(s, faulty + 1, total - (faulty + 1))
	fmt.Println("Shards are ", shards)
	//Get the string version of the string
	hashStr := ConvertBytesToString(Hash([]byte(s)))

	for i := 0; i < total; i++ {
		code := ConvertBytesToString(shards[i])
		m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, HashData: hashStr, Data: code, Round:1}
		sendReq := PrepareSend{M: m, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func SimpleECMessageHandler(m Message) {
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

		SimpleECCheck(m)
	}
}

func SimpleECEchoHandler(m Message) {
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
				expectedHash := echo.GetHashData()
				inputHash := ConvertBytesToString(Hash([]byte(v)))

				if compareHash(expectedHash, inputHash) {
					fmt.Println("Include ", v)
					DataSet[v] = echo.GetHashData()
					break
				}
			}
		}
		SimpleECCheck(m)
	}
}

func SimpleECCheck(m Message) {
	dataExist, data := checkDataExist(m.GetHashData())

	if dataExist {
		echo := ECHOStruct{Header:ECHO, Id:m.GetId(), HashData:m.GetHashData(), Round: m.GetRound()}
		if EchoRecCountSet[echo] >= faulty + 1 {
			identifier := m.GetId() + strconv.Itoa(m.GetRound())
			if _, sent := EchoSentSet[identifier]; !sent {
				EchoSentSet[identifier] = true
				shards := Encode(data, faulty + 1, total - (faulty + 1))

				code := ConvertBytesToString(shards[serverMap[MyID]])
				echo = ECHOStruct{Header:ECHO, Id:m.GetId(), HashData:m.GetHashData(), Round: m.GetRound(), Data: code, SenderId:MyID}
				sendReq := PrepareSend{M: echo, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}
		if EchoRecCountSet[echo] >= total - faulty {
			if _, e := acceptData[data]; ! e {
				acceptData[data] = true
				fmt.Println("Reliable Accept " + data)
			}
		}
	}
}


func compareHash(expectedHash, inputHash string) bool{
	return expectedHash == inputHash
}