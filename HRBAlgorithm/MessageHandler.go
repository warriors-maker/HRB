package HRBAlgorithm

import (
	"fmt"
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
		fmt.Printf("%+v\n",m)

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
	fmt.Printf("Req: %+v\n",data)
	if _, ok := ReqReceiveSet[identifier]; !ok {
		ReqReceiveSet[identifier] = true
		if exist, d := checkDataExist(data.GetHashData()); exist {
			//send val to the requested id
			fwdSendMsg := FWDStruct{Header:FWD, Id: data.GetId(), Round:data.GetRound(), SenderId:MyID, Data: d}
			req := PrepareSend{M: fwdSendMsg, SendTo: data.GetSenderId()}
			SendReqChan <- req
			return true, true
		}
		return true, false
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
	fmt.Printf("Fwd: %+v\n",data)
	m := REQStruct{Header:REQ, Id:data.GetId(), HashData:hashStr, Round: data.GetRound(), SenderId:MyID}
	if hasSent(ReqSentSet[m], data.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			FwdReceiveSet[identifier] = true
			DataSet[data.GetData()] = hashStr
			//check
			check(m)
			return true, true
		}
		return true, false
	}
	return false, false

}

func EchoHandler (data Message) (bool, int){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
		fmt.Printf("%+v\n",m)
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

		if len(AccRecCountSet[accM]) == faulty + 1 {
			//Same set
			ReqSentSet[reqMes] = AccRecCountSet[accM]
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				//Send Req to these f + 1
				for _,id := range ReqSentSet[reqMes] {
					reqMessage := PrepareSend{M: reqMes, SendTo: id}
					SendReqChan <- reqMessage
				}
			}
			return true, len(AccRecCountSet[accM]), true
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
		return true, len(AccRecCountSet[accM]), false
	}
	return false, 0, false
}

func check(m Message) []bool {
	fmt.Println("Inside check")
	flags := []bool{false, false, false, false}
	if exist, value := checkDataExist(m.GetHashData()); exist {
		echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
		acc := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}

		fmt.Printf("%+v\n",echo)
		fmt.Printf("%+v\n",acc)

		identifier := m.GetId() + strconv.Itoa(m.GetRound())

		if EchoRecCountSet[echo] >= faulty + 1 {
			fmt.Println("Receive more than faulty + 1 echo message")
			if _, sent := EchoSentSet[identifier]; !sent {
				fmt.Println("Sent Echo to other servers")
				echo.SetSenderId(MyID)
				EchoSentSet[identifier] = true
				sendReq := PrepareSend{M: echo, SendTo:"all"}
				SendReqChan <- sendReq
			}
		}

		if EchoRecCountSet[echo] >= total - faulty {
			fmt.Println("Receive more than total - faulty echo message")
			if _, sent := AccSentSet[identifier]; !sent {
				AccSentSet[identifier] = true
				acc.SetSenderId(MyID)
				//send ACC to all servers
				fmt.Println("Sent Acc to other servers")
				sendReq := PrepareSend{M: acc, SendTo:"all"}
				SendReqChan <- sendReq
				flags[1] = true
			}
		}

		if len(AccRecCountSet[acc]) >= faulty + 1 {
			fmt.Println("Identifier"+identifier)
			if _,sent := AccSentSet[identifier]; !sent {
				AccSentSet[identifier] = true
				acc.SetSenderId(MyID)
				//send ACC to all servers
				fmt.Println("Sent Acc to other servers")
				sendReq := PrepareSend{M: acc, SendTo:"all"}
				SendReqChan <- sendReq
				flags[2] = true
			}
		}

		if len(AccRecCountSet[acc]) >= total - faulty {
			//accept
			fmt.Println("Reliable Accept " + value)
			flags[3] = true
		}
	}
	return flags
}