package Server

import (
	"HRB/HRBAlgorithm"
	"encoding/gob"
	"net"
	"time"
)

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
			//if crashFailure, just do not send the data
			if algorithm == 7 || algorithm == 9{
				if source {
					if ipPort == serverList[1] {

					} else {
						encoder.Encode(&data)
					}
				} else {

				}
			} else if algorithm ==5 || algorithm == 6 {
				if source {
					if data.Message.GetHeaderType() == HRBAlgorithm.MSG {
						if ipPort == serverList[1] {
							wrongMessage := HRBAlgorithm.MSGStruct{Header:HRBAlgorithm.MSG, Id:data.Message.GetId(), SenderId:data.Message.GetSenderId(), Data:"Asdaw!2heyhe", Round:data.Message.GetRound()}
							wrongTcpMessage := TcpMessage{Message:wrongMessage}
							encoder.Encode(&wrongTcpMessage)
						} else {
							encoder.Encode(&data)
						}
					} else {
						encoder.Encode(&data)
					}
				}
			} else {
				//Source Equivocate
				if source {
					if data.Message.GetHeaderType() == HRBAlgorithm.MSG {
						if ipPort == serverList[1] {
							//fmt.Println("Do not send to" + ipPort)
						} else {
							encoder.Encode(&data)
						}
					} else {
						encoder.Encode(&data)
					}
				} else if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
					//fmt.Println("Set data to null")
					correct := data.Message
					// Create a Faulty Message

					faulty := HRBAlgorithm.ECHOStruct{Id:correct.GetId(), Data: data.Message.GetData()+ "1asdadadwa", SenderId:correct.GetSenderId(),
						HashData:")a1s2f*(", Round:correct.GetRound(), Header:HRBAlgorithm.ECHO}
					data = TcpMessage{Message:faulty}

					encoder.Encode(&data)
				} else {
					if algorithm == 5 || algorithm == 6 {
						encoder.Encode(&data)
					} else if counter % 2 == 0 {
						encoder.Encode(&data)
					}
				}
			}
		} else {
			//fmt.Printf("Benchmark Send Data Externally to %+v\n",data)
			encoder.Encode(&data)
		}
	}
}



