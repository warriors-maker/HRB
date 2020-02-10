package Server

import (
"HRB/HRBAlgorithm"
"encoding/gob"
"fmt"
"net"
"time"
)

//ipPort: the targer ipAddress to write to
func ExternalTcpWriter(ipPort string, ch chan TcpMessage) {
	fmt.Println("Benchmark Channel for sending data to " + ipPort)

	conn, err:= net.Dial("tcp",ipPort)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp",ipPort)
		time.Sleep(10*time.Second)
	}

	encoder := gob.NewEncoder(conn)
	for {
		data := <-ch
		if isFault {
			if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
				fmt.Println("Set data to null")
				correct := data.Message
				// Create a Faulty Message

				faulty := HRBAlgorithm.ECHOStruct{Id:correct.GetId(), Data: data.Message.GetData()+ "1asdadadwa", SenderId:correct.GetSenderId(),
					HashData:")a1s2f*(", Round:correct.GetRound(), Header:HRBAlgorithm.ECHO}
				data = TcpMessage{Message:faulty}

				encoder.Encode(&data)
			} else {
				encoder.Encode(&data)
			}
		} else {
			fmt.Printf("Send Data to %+v\n",data)
			encoder.Encode(&data)
		}
	}
}



