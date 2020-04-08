package HRBAlgorithm

import (
	"HRB/HRBMessage"
	"fmt"
	"strconv"
	"time"
)

func SimpleBroadcast(byteLength, round int) {
	time.Sleep(3 *time.Second)
	for i := 0; i < round; i++ {
		//if i % 200 == 0 {
		//	time.Sleep(1*time.Second)
		//}
		s := RandStringBytes(byteLength)
		m := HRBMessage.MSGStruct{Id: MyID, SenderId:MyID, Data: s, Header:HRBMessage.MSG, Round:i}
		for _, server := range serverList {
			//fmt.Println("Protocal send to ", server)
			sendReq := HRBMessage.PrepareSend{M: m, SendTo: server}
			SendReqChan <- sendReq
		}
	}
}


func Msghandler(d HRBMessage.Message) (bool, int, string){
	data,_ := d.(HRBMessage.MSGStruct)


	identifier := data.GetId() + strconv.Itoa(data.GetRound())
	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true


		//include the data with key the original data and val its hashvalue
		hashstr := ConvertBytesToString(Hash([]byte(data.GetData())))


		//Main logic
		m := HRBMessage.ECHOStruct{Header: HRBMessage.ECHO, Id:data.GetId(), HashData:hashstr, Round: data.GetRound(), SenderId:MyID}
		//fmt.Printf("%+v\n",m)

		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
			sendReq := HRBMessage.PrepareSend{M: m, SendTo:"all"}
			SendReqChan <- sendReq
		}
		//Check
		check(m)
		return true, EchoRecCountSet[m],hashstr
	}
	return false, 0,""
}

func ReqHandler (data HRBMessage.Message) (bool, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	//fmt.Printf("Req: %+v\n",data)
	if _, ok := ReqReceiveSet[identifier]; !ok {
		ReqReceiveSet[identifier] = true
		if exist, d := checkDataExist(data.GetHashData()); exist {
			//send val to the requested id
			fwdSendMsg := HRBMessage.FWDStruct{Header: HRBMessage.FWD, Id: data.GetId(), Round:data.GetRound(), SenderId:MyID, Data: d}
			req := HRBMessage.PrepareSend{M: fwdSendMsg, SendTo: data.GetSenderId()}
			SendReqChan <- req
			//return true, true
		}
		//return true, false
	}
	return false, false
}

//Need to check this function later
func FwdHandler (data HRBMessage.Message) (bool, bool){
	identifier := data.GetId() + strconv.Itoa(data.GetRound()) + data.GetSenderId();
	/*
		have seen echo + 1
	*/
	//data type check
	hashStr := ConvertBytesToString(Hash([]byte(data.GetData())))
	//fmt.Printf("Fwd: %+v\n",data)
	m := HRBMessage.REQStruct{Header: HRBMessage.REQ, Id:data.GetId(), HashData:hashStr, Round: data.GetRound(), SenderId:MyID}
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

func EchoHandler (data HRBMessage.Message) (bool, int){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true
		DataSet[data.GetData()] = data.GetHashData()
		m := HRBMessage.ECHOStruct{Header: HRBMessage.ECHO, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
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


func AccHandler (data HRBMessage.Message) (bool, int, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	//fmt.Println("DataID "+data.GetId())
	//fmt.Println("AccHandler identifier " + identifier)

	if _, seen := AccReceiveSet[identifier]; !seen {
		AccReceiveSet[identifier] = true

		accM := HRBMessage.ACCStruct{Header: HRBMessage.ACC, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
		if l, ok := AccRecCountSet[accM]; !ok {
			l = []string{data.GetSenderId()}
			AccRecCountSet[accM] = l
		} else {
			l = append(l, data.GetSenderId())
			AccRecCountSet[accM] = l
		}

		reqMes := HRBMessage.REQStruct{Header: HRBMessage.REQ, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound(), SenderId:MyID}
		//mt.Printf("Acc in AccHandler %+v, %+v \n", accM, AccRecCountSet[accM])
		if len(AccRecCountSet[accM]) == faulty + 1 {
			//Same set
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				ReqSentSet[reqMes] = AccRecCountSet[accM]
				//Send Req to these f + 1
				for _,id := range ReqSentSet[reqMes] {
					reqMessage := HRBMessage.PrepareSend{M: reqMes, SendTo: id}
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
				reqMessage := HRBMessage.PrepareSend{M: reqMes, SendTo: data.GetSenderId()}
				SendReqChan <- reqMessage
				return true, len(AccRecCountSet[accM]), true
			}
		}
		check(data)

	}
	return false, 0, false
}

func check(m HRBMessage.Message) []bool {
	//fmt.Println(m.GetHashData())
	flags := []bool{false, false, false, false}

	if exist, _ := checkDataExist(m.GetHashData()); exist {
		echo := HRBMessage.ECHOStruct{Header: HRBMessage.ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
		acc := HRBMessage.ACCStruct{Header: HRBMessage.ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}


		identifier := m.GetId() + strconv.Itoa(m.GetRound())

		if EchoRecCountSet[echo] >= faulty + 1 {
			if _, sent := EchoSentSet[identifier]; !sent {
				//fmt.Println("Receive more than faulty + 1 echo message")
				//fmt.Println("Sent Echo to other servers")
				echoSend := HRBMessage.ECHOStruct{Header: HRBMessage.ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				EchoSentSet[identifier] = true
				sendReq := HRBMessage.PrepareSend{M: echoSend, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		if EchoRecCountSet[echo] >= total - faulty {
			if _, sent := AccSentSet[identifier]; !sent {
				//fmt.Println("Receive more than total - faulty echo message")
				AccSentSet[identifier] = true

				accSend := HRBMessage.ACCStruct{Header: HRBMessage.ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				//fmt.Println("Acc sender ID "+ MyID + "," + acc.SenderId)
				//fmt.Printf("Sent Acc to other servers %+v \n" , acc)
				sendReq := HRBMessage.PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
				flags[1] = true
			}
		}

		//fmt.Println("Acc Info",len(AccRecCountSet[acc]), AccRecCountSet[acc], faulty + 1)

		if len(AccRecCountSet[acc]) >= faulty + 1 {
			if _,sent := AccSentSet[identifier]; !sent {
				//fmt.Println("Receive more than f + 1 Acc message")
				AccSentSet[identifier] = true
				accSend := HRBMessage.ACCStruct{Header: HRBMessage.ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				//send ACC to all servers
				sendReq := HRBMessage.PrepareSend{M: accSend, SendTo:"all"}
				SendReqChan <- sendReq
				flags[2] = true
			}
		}

		fmt.Println(AccRecCountSet[acc], total, faulty)

		if len(AccRecCountSet[acc]) >= total - faulty {
			if _, e :=acceptData[identifier]; !e {
				acceptData[identifier] = true
				if algorithm == 8 {
					for data, hash := range DataSet {
						if hash == m.GetHashData() {
							hashTag := HRBMessage.Binary{Header: HRBMessage.BIN, Round:m.GetRound(), Id:m.GetId(), HashData:data, SenderId:MyID}
							fmt.Printf("ReliableAccept %+v \n" , hashTag)
							sendReq := HRBMessage.PrepareSend{M: hashTag, SendTo: MyID}
							SendReqChan <- sendReq
							break
						}
					}
				} else {
					stats := HRBMessage.StatStruct{Id: m.GetId(), Round: m.GetRound(), Header: HRBMessage.Stat}
					statInfo := HRBMessage.PrepareSend{M: stats, SendTo:MyID}
					SendReqChan <- statInfo
				}
			}
		}
	}
	return flags
}

/*
Helper Function
*/

func hasSent(l []string, val string) bool{
	for _, v := range l {
		if v == val {
			return true
		}
	}
	return false
}

func checkDataExist(expectedHash string) (bool, string) {
	for k,v := range DataSet {
		if v == expectedHash {
			return true, k
		}
	}
	return false,""
}
