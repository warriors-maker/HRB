package HRBAlgorithm

import (
	"strconv"
)


func Msghandler (d Message) (bool, int, string){
	data,_ := d.(MSGStruct)


	identifier := data.GetId() + strconv.Itoa(data.GetRound())
	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true


		//include the data with key the original data and val its hashvalue
		hashstr := ConvertBytesToString(Hash([]byte(data.GetData())))

		DataSet[data.GetData()] = hashstr

		//var hash []byte
		//var hashstr string

		/*
			have seen echo + 1
		*/

		//data type check
		//oknum,num := Util.ParseInt(data.GetData())
		//if oknum {
		//	hash = Hash([]byte(strconv.Itoa(num)))
		//	hashstr = ConvertBytesToString(hash)
		//}
		//
		//okstr ,str := Util.ParseString(data.GetData())
		//if okstr {
		//	hash = Hash([]byte(str))
		//	hashstr = ConvertBytesToString(hash)
		//}

		//Main logic
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:hashstr, Round: data.GetRound(), SenderId:MyID}
		//fmt.Printf("%+v\n",m)

		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
			sendReq := PrepareSend{M: m, SendTo:"all"}
			SendReqChan <- sendReq
		}
		//Check
		check(m)
		return true, EchoRecCountSet[m],hashstr
	}
	return false, 0,""
}

func ReqHandler (data Message) (bool, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	//fmt.Printf("Req: %+v\n",data)
	if _, ok := ReqReceiveSet[identifier]; !ok {
		ReqReceiveSet[identifier] = true
		if exist, d := checkDataExist(data.GetHashData()); exist {
			//send val to the requested id
			fwdSendMsg := FWDStruct{Header:FWD, Id: data.GetId(), Round:data.GetRound(), SenderId:MyID, Data: d}
			req := PrepareSend{M: fwdSendMsg, SendTo: data.GetSenderId()}
			SendReqChan <- req
			//return true, true
		}
		//return true, false
	}
	return false, false
}

//Need to check this function later
func FwdHandler (data Message) (bool, bool){
	identifier := data.GetId() + strconv.Itoa(data.GetRound()) + data.GetSenderId();
	/*
		have seen echo + 1
	*/
	//data type check
	hashStr := ConvertBytesToString(Hash([]byte(data.GetData())))
	//fmt.Printf("Fwd: %+v\n",data)
	m := REQStruct{Header:REQ, Id:data.GetId(), HashData:hashStr, Round: data.GetRound(), SenderId:MyID}
	if hasSent(ReqSentSet[m], data.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			FwdReceiveSet[identifier] = true
			DataSet[data.GetData()] = hashStr
			//check
			check(m)
			//return true, true
		}
		//return true, false
	}
	return false, false

}

func EchoHandler (data Message) (bool, int){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
		//fmt.Printf("%+v\n",m)
		if count, ok := EchoRecCountSet[m]; ok {
			EchoRecCountSet[m] = count + 1
		} else {
			EchoRecCountSet[m] = 1
		}

		//Check
		check(m)
		//flags := check(data)
		return true, EchoRecCountSet[m]
	}
	return false,0
}


func AccHandler (data Message) (bool, int, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	//fmt.Println("DataID "+data.GetId())
	//fmt.Println("AccHandler identifier " + identifier)

	if _, seen := AccReceiveSet[identifier]; !seen {
		AccReceiveSet[identifier] = true

		accM := ACCStruct{Header: ACC, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
		if l, ok := AccRecCountSet[accM]; !ok {
			l = []string{data.GetSenderId()}
			AccRecCountSet[accM] = l
		} else {
			l = append(l, data.GetSenderId())
			AccRecCountSet[accM] = l
		}

		reqMes := REQStruct{Header:REQ, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound(), SenderId:MyID}
		//mt.Printf("Acc in AccHandler %+v, %+v \n", accM, AccRecCountSet[accM])
		if len(AccRecCountSet[accM]) == faulty + 1 {
			//Same set
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				ReqSentSet[reqMes] = AccRecCountSet[accM]
				//Send Req to these f + 1
				for _,id := range ReqSentSet[reqMes] {
					reqMessage := PrepareSend{M: reqMes, SendTo: id}
					SendReqChan <- reqMessage
				}
			}

		} else if len(AccRecCountSet[accM]) > faulty + 1 {
			//send the request to the current id
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				//Todo: Need to set a flag that if we have accepted this round, we donot need to send the req again
				l := ReqSentSet[reqMes]
				l = append(l, data.GetSenderId())
				ReqSentSet[reqMes] = l
				//send request
				reqMessage := PrepareSend{M: reqMes, SendTo: data.GetSenderId()}
				SendReqChan <- reqMessage
				return true, len(AccRecCountSet[accM]), true
			}
		}
		check(data)

	}
	return false, 0, false
}

func check(m Message) []bool {
	//fmt.Println("Inside check")
	//fmt.Println(m.GetHashData())
	flags := []bool{false, false, false, false}

	if exist, _ := checkDataExist(m.GetHashData()); exist {
		echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
		acc := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}


		identifier := m.GetId() + strconv.Itoa(m.GetRound())

		if EchoRecCountSet[echo] >= faulty + 1 {
			if _, sent := EchoSentSet[identifier]; !sent {
				//fmt.Println("Receive more than faulty + 1 echo message")
				//fmt.Println("Sent Echo to other servers")
				echoSend := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				EchoSentSet[identifier] = true
				sendReq := PrepareSend{M: echoSend, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		if EchoRecCountSet[echo] >= total - faulty {
			if _, sent := AccSentSet[identifier]; !sent {
				//fmt.Println("Receive more than total - faulty echo message")
				AccSentSet[identifier] = true

				accSend := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				//fmt.Println("Acc sender ID "+ MyID + "," + acc.SenderId)
				//fmt.Printf("Sent Acc to other servers %+v \n" , acc)
				sendReq := PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
				flags[1] = true
			}
		}

		//fmt.Println("Acc Info",len(AccRecCountSet[acc]), AccRecCountSet[acc], faulty + 1)

		if len(AccRecCountSet[acc]) >= faulty + 1 {
			if _,sent := AccSentSet[identifier]; !sent {
				//fmt.Println("Receive more than f + 1 Acc message")
				AccSentSet[identifier] = true
				accSend := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				sendReq := PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
				flags[2] = true
			}
		}

		if len(AccRecCountSet[acc]) >= total - faulty {
			if _, e :=acceptData[identifier]; !e {
				acceptData[identifier] = true
				if algorithm == 8 {
					for data, hash := range DataSet {
						if hash == m.GetHashData() {
							hashTag := Binary{Header:BIN, Round:m.GetRound(), Id:m.GetId(), HashData:data, SenderId:MyID}
							//fmt.Printf("ReliableAccept %+v \n" , hashTag)
							sendReq := PrepareSend{M: hashTag, SendTo: MyID}
							SendReqChan <- sendReq
							break
						}
					}
				} else {
					stats := StatStruct{Id:m.GetId(), Round: m.GetRound(), Header:Stat}
					statInfo :=PrepareSend{M:stats, SendTo:MyID}
					SendReqChan <- statInfo
				}
			}
		}
	}
	return flags
}

