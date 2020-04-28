package Server

import (
	"HRB/HRBMessage"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

/*
Write to external nodes
 */

//ipPort: the targer ipAddress to write to
func ExternalTcpWriter(ipPort string, ch chan TcpMessage) {
	//fmt.Println("Benchmark Channel for sending data to " + ipPort)

	conn, err:= net.Dial("tcp",ipPort)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp",ipPort)
		time.Sleep(2*time.Second)
	}

	encoder := gob.NewEncoder(conn)
	for {
		counter := 0
		data := <-ch
		counter = counter + 1
		if isFault {
			switch v := data.Message.(type) {
			case HRBMessage.MSGStruct:
				if ipPort == serverList[1] {
					//fmt.Println("Do not send to" + ipPort)
				} else {
					encoder.Encode(&data)
				}
				break
			case HRBMessage.ECHOStruct:
				correct := v
				// Create a Faulty Message

				faulty := HRBMessage.ECHOStruct{Id: correct.GetId(), Data: correct.GetData()+ "1asdadadwa", SenderId:correct.GetSenderId(),
					HashData:")a1s2f*(", Round:correct.GetRound(), Header: HRBMessage.ECHO}
				data = TcpMessage{Message:faulty}

				encoder.Encode(&data)
				break
			case HRBMessage.ACCStruct:
				encoder.Encode(&data)
				//fmt.Println("Acc")
				break
			case HRBMessage.REQStruct:
				encoder.Encode(&data)
				//fmt.Println("Req")
				break
			case HRBMessage.FWDStruct:
				encoder.Encode(&data)
				//fmt.Print("FWD")
				break
			default:
				fmt.Printf("Sending : %+v\n", v)
				fmt.Println("I do ot understand what you send")
			}

			} else {
			//fmt.Printf("Benchmark Send Data Externally to %+v\n",data)
			encoder.Encode(&data)
		}
	}
}



