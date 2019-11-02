package HRBAlgorithm

import (
	"fmt"
	"strconv"
)

//Same as the more complicated one
func SimpleMsgHandler(d Message) {
	data,ok := d.(MSGStruct)
	if ok {
		fmt.Println(ok)
	}

	identifier := data.GetId() + strconv.Itoa(data.GetRound())
	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true

		hashStr := ConvertBytesToString([]byte(data.GetData()))

		//include the data with key the original data and val its hashvalue
		DataSet[data.GetData()] = hashStr

		//Main logic
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:hashStr, Round: data.GetRound()}
		fmt.Printf("%+v\n",m)

		if l, ok := simpleEchoRecCountSet[m]; ok {
			l = append(l, d.GetSenderId())
			simpleEchoRecCountSet[m] = l
		} else {
			simpleEchoRecCountSet[m] = []string{d.GetSenderId()}
		}

		//ToDo: Send Echo to all servers
		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
		}

	}

}

func SimpleReqHandler(d Message) {
	ReqHandler(d)
}

func SimpleFwdHandler(d Message) {
	FwdHandler(d)
}

func SimpleEchoHandler(d Message) {
	identifier := d.GetId()+ strconv.Itoa(d.GetRound()) + d.GetSenderId()
	if _,seen := EchoReceiveSet[identifier]; !seen {
		EchoReceiveSet[identifier] = true

		//update count
		m := ECHOStruct{Header:ECHO, Id:d.GetId(), HashData:d.GetHashData(), Round: d.GetRound()}
		fmt.Printf("%+v\n",m)
		if l, ok := simpleEchoRecCountSet[m]; ok {
			l := append(l, m.GetSenderId())
			simpleEchoRecCountSet[m] = l
		} else {
			l := []string{d.GetSenderId()}
			simpleEchoRecCountSet[m] = l
		}

		reqMes := REQStruct{Header:REQ, Id:d.GetId(), HashData:d.GetHashData(), Round: d.GetRound()}
		if len(simpleEchoRecCountSet[m]) == faulty + 1 {
			ReqSentSet[reqMes] = simpleEchoRecCountSet[m]
			if exist, _ := checkDataExist(d.GetHashData()); !exist {
				//Todo: Send Req to these f + 1
				for _,id := range simpleEchoRecCountSet[m] {
					fmt.Println(id)
				}
			}

		} else if len(simpleEchoRecCountSet[m]) > faulty + 1 {
			if exist, _ := checkDataExist(d.GetHashData()); !exist {
				//send the request to the current id
				//Todo: Need to set a flag that if we have accepted this round, we do not need to send the req again
				l := ReqSentSet[reqMes]
				l = append(l, d.GetSenderId())
				ReqSentSet[reqMes] = l
				//send request
			}
		}

		//Check
		//flags := check(data)

	}

}

func SimpleCheck(m Message) {
	fmt.Println("Inside check")
	echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}

	fmt.Printf("%+v\n",echo)

	identifier := m.GetId() + strconv.Itoa(m.GetRound())
	flags := []bool{false, false}

	if EchoRecCountSet[echo] >= faulty + 1 {
		fmt.Println("Receive more than faulty + 1 echo message")
		if _, sent := EchoSentSet[identifier]; !sent {
			fmt.Println("Sent Echo to other servers")
			EchoSentSet[identifier] = true
			//send Echo to all servers
			flags[0] = true
		}
	}

	if EchoRecCountSet[echo] >= total - faulty {
		fmt.Println("Receive more than total - faulty echo message")
		if _, sent := AccSentSet[identifier]; !sent {
			AccSentSet[identifier] = true
			fmt.Println("Simple Reliable Accept")
			flags[1] = true
		}
	}


}