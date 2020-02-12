package HRBAlgorithm

import (
	"strconv"
	"time"
)

func SimpleBroadcast(byteLength, round int) {
	time.Sleep(3 *time.Second)
	for i := 0; i < round; i++ {
		s := RandStringBytes(byteLength)
		m := MSGStruct{Id: MyID, SenderId:MyID, Data: s, Header:MSG, Round:i}
		for _, server := range serverList {
			//fmt.Println("Protocal send to ", id)
			sendReq := PrepareSend{M: m, SendTo: server}
			SendReqChan <- sendReq
		}
	}
}

//Same as the more complicated one
func SimpleMsgHandler(d Message) {
	data,_ := d.(MSGStruct)

	identifier := data.GetId() + strconv.Itoa(data.GetRound())
	if _, seen := MessageReceiveSet[identifier]; !seen {

		MessageReceiveSet[identifier] = true

		hashStr := ConvertBytesToString(Hash([]byte(data.GetData())))

		//include the data with key the original data and val its hashvalue
		DataSet[data.GetData()] = hashStr

		//Main logic
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:hashStr, Round: data.GetRound(), SenderId:MyID}

		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
			sendReq := PrepareSend{M: m, SendTo:"all"}
			SendReqChan <- sendReq
		}
		SimpleCheck(d)
	}

}

func SimpleReqHandler(d Message) {
	ReqHandler(d)
}

func SimpleFwdHandler(data Message) {
	identifier := data.GetId() + strconv.Itoa(data.GetRound()) + data.GetSenderId();
	/*
		have seen echo + 1
	*/
	//data type check
	hashStr := ConvertBytesToString(Hash([]byte(data.GetData())))
	//fmt.Printf("Fwd: %+v\n",data)
	m := REQStruct{Header:REQ, Id:data.GetId(), HashData:hashStr, Round: data.GetRound(), SenderId:MyID}
	//fmt.Printf("ReceiveBack FWD %+v\n" , m)
	if hasSent(ReqSentSet[m], data.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			//fmt.Println("Receive fwd back from the request")
			FwdReceiveSet[identifier] = true
			DataSet[data.GetData()] = hashStr
			//check
			SimpleCheck(m)
		}
	}
}

func SimpleEchoHandler(d Message) {
	identifier := d.GetId()+ strconv.Itoa(d.GetRound()) + d.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		//update count
		m := ECHOStruct{Header:ECHO, Id:d.GetId(), HashData:d.GetHashData(), Round: d.GetRound()}
		//fmt.Printf("Echo: %+v\n",m)
		if l, ok := simpleEchoRecCountSet[m]; ok {
			l := append(l, d.GetSenderId())
			simpleEchoRecCountSet[m] = l
		} else {
			l := []string{d.GetSenderId()}
			simpleEchoRecCountSet[m] = l
		}

		reqMes := REQStruct{Header:REQ, Id:d.GetId(), HashData:d.GetHashData(), Round: d.GetRound(), SenderId:MyID}
		if len(simpleEchoRecCountSet[m]) == faulty + 1 {
			ReqSentSet[reqMes] = simpleEchoRecCountSet[m]

			if exist, _ := checkDataExist(d.GetHashData()); !exist {
				//Todo: Send Req to these f + 1
				//fmt.Printf("Send request to f + 1 servers: %+v\n", reqMes)
				for _,id := range simpleEchoRecCountSet[m] {
					reqMessage := PrepareSend{M: reqMes, SendTo: id}
					SendReqChan <- reqMessage
				}
			}
			SimpleCheck(d)
		} else if len(simpleEchoRecCountSet[m]) > faulty + 1 {
			if exist, _ := checkDataExist(d.GetHashData()); !exist {
				//send the request to the current id
				//Todo: Need to set a flag that if we have accepted this round, we do not need to send the req again
				//fmt.Println("Send request to individual server: ", ReqSentSet[reqMes])
				l := ReqSentSet[reqMes]
				l = append(l, d.GetSenderId())
				ReqSentSet[reqMes] = l
				//send request
				reqMessage := PrepareSend{M: reqMes, SendTo: d.GetSenderId()}
				SendReqChan <- reqMessage
			}
			SimpleCheck(d)
		}
	}

}

func SimpleCheck(m Message) {
	if exist, value := checkDataExist(m.GetHashData()); exist {
		echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
		//fmt.Printf("%+v\n",echo)

		identifier := m.GetId() + strconv.Itoa(m.GetRound())
		flags := []bool{false, false}

		//fmt.Println(simpleEchoRecCountSet[echo], len(simpleEchoRecCountSet[echo]))
		if len(simpleEchoRecCountSet[echo]) >= faulty + 1 {
			//fmt.Println("Receive more than faulty + 1 echo message")
			if _, sent := EchoSentSet[identifier]; !sent {
				//fmt.Println("Sent Echo to all servers")
				EchoSentSet[identifier] = true
				//send Echo to all servers
				echoSend := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId(), SenderId:MyID}
				req := PrepareSend{M:echoSend, SendTo:"all"}
				SendReqChan <- req
				flags[0] = true
			}
		}

		if len(simpleEchoRecCountSet[echo]) >= total - faulty {
			//fmt.Println("Receive more than total - faulty echo message")
			if _, e := acceptData[value]; ! e {
				acceptData[value] = true
				stats := StatStruct{Id:m.GetId(), Round: m.GetRound(), Header:Stat}
				statInfo :=PrepareSend{M:stats, SendTo:MyID}
				SendReqChan <- statInfo
			}
		}

		//fmt.Printf("Count: %d\n",len(simpleEchoRecCountSet[echo]))
	}

}