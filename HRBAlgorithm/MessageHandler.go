package HRBAlgorithm

import (
	"HRB/Util"
	"fmt"
	"strconv"
)


func Msghandler (d Message) (bool, int, string){
	data,ok := d.(MSGStruct)
	if ok {
		fmt.Println(ok)
	}

	identifier := data.GetId() + strconv.Itoa(data.GetRound())
	if _, seen := MessageReceiveSet[identifier]; !seen {
		MessageReceiveSet[identifier] = true

		//include the data with key the original data and val its hashvalue
		DataSet[data.GetData()] = ConvertBytesToString(Hash([]byte(data.GetData())))

		var hash []byte
		var hashstr string

		/*
			have seen echo + 1
		*/

		//data type check
		oknum,num := Util.ParseInt(data.GetData())
		if oknum {
			hash = Hash([]byte(strconv.Itoa(num)))
			hashstr = ConvertBytesToString(hash)
		}

		okstr ,str := Util.ParseString(data.GetData())
		if okstr {
			hash = Hash([]byte(str))
			hashstr = ConvertBytesToString(hash)
		}

		//Main logic
		m := ECHOStruct{Header:ECHO, Id:data.GetId(), HashData:hashstr, Round: data.GetRound()}
		fmt.Printf("%+v\n",m)

		if count, ok := EchoRecCountSet[m]; ok {
			EchoRecCountSet[m] = count + 1
		} else {
			EchoRecCountSet[m] = 1
		}

		// Send Echo to all servers
		if _, sent := EchoSentSet[identifier]; !sent {
			EchoSentSet[identifier] = true
		}
		return true, EchoRecCountSet[m],hashstr
	}
	return false, 0,""
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
		//flags := check(data)
		return true, EchoRecCountSet[m]
	}
	return false,0
}

func AccHandler (data Message) (bool, int, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	if _, seen := AccReceiveSet[identifier]; !seen {
		AccReceiveSet[identifier] = true

		m := ACCStruct{Header:ACC, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}
		reqMes := REQStruct{Header:REQ, Id:data.GetId(), HashData:data.GetHashData(), Round: data.GetRound()}

		if l, ok := AccRecCountSet[m]; !ok {
			l = []string{data.GetSenderId()}
			AccRecCountSet[m] = l
		} else {
			l = append(l, data.GetSenderId())
			AccRecCountSet[m] = l
		}

		if len(AccRecCountSet[m]) == faulty + 1 {
			//Same set
			ReqSentSet[reqMes] = AccRecCountSet[m]
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				//Send Req to these f + 1
				for _,id := range AccRecCountSet[m] {
					fmt.Println(id)
				}
			}
			return true, len(AccRecCountSet[m]), true
		} else if len(AccRecCountSet[m]) > faulty + 1 {
			//send the request to the current id
			if exist, _ := checkDataExist(data.GetHashData()); !exist {
				//Todo: Need to set a flag that if we have accepted this round, we donot need to send the req again
				l := ReqSentSet[reqMes]
				l = append(l, data.GetSenderId())
				ReqSentSet[reqMes] = l
				//send request
				return true, len(AccRecCountSet[m]), true
			}
		}

		return true, len(AccRecCountSet[m]), false
	}
	return false, 0, false
}

func ReqHandler (data Message) (bool, bool){
	identifier := data.GetId()+ strconv.Itoa(data.GetRound()) + data.GetSenderId()
	fmt.Printf("Req: %+v\n",data)
	if _, ok := ReqReceiveSet[identifier]; !ok {
		ReqReceiveSet[identifier] = true
		if exist, _ := checkDataExist(data.GetHashData()); exist {
			//send val to the requested id
			//fmt.Println(val, data.GetHashData())
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
	m := REQStruct{Header:REQ, Id:data.GetId(), HashData:hashStr, Round: data.GetRound()}
	if hasSent(ReqSentSet[m], data.GetSenderId()) {
		if _,seen := FwdReceiveSet[identifier]; !seen {
			FwdReceiveSet[identifier] = true
			DataSet[data.GetData()] = hashStr
			//check
			return true, true
		}
		return true, false
	}
	return false, false

}

func check(m Message) []bool {
	fmt.Println("Inside check")
	echo := ECHOStruct{Header:ECHO, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}
	acc := ACCStruct{Header:ACC, Round:m.GetRound(), HashData:m.GetHashData(), Id:m.GetId()}

	fmt.Printf("%+v\n",echo)
	fmt.Printf("%+v\n",acc)

	identifier := m.GetId() + strconv.Itoa(m.GetRound())

	flags := []bool{false, false, false, false}

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
			//send ACC to all servers
			fmt.Println("Sent Acc to other servers")
			flags[1] = true
		}
	}

	if len(AccRecCountSet[acc]) >= faulty + 1 {
		fmt.Println("Identifier"+identifier)
		if _,sent := AccSentSet[identifier]; !sent {
			AccSentSet[identifier] = true
			//send ACC to all Servers
			flags[2] = true
		}
	}

	if len(AccRecCountSet[acc]) >= total - faulty {
		//accept
			flags[3] = true
	}
	return flags
}