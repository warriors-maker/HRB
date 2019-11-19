package Server

import (
	"HRB/HRBAlgorithm"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

//ipPort: the targer ipAddress to write to
func TcpWriter(ipPort string, ch chan TcpMessage) {
	counter := 0
	fmt.Println("Channel for sending data to " + ipPort)

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

			//if counter % 2 == 0 {
			//
			//} else {
			//	fmt.Printf("Sending : %+v to %v\n", data.Message, ipPort)
			//	//time.Sleep(7*time.Second)
			//	encoder.Encode(&data)
			//}

			/*
			Test the Extreme Case for EC Coding
			 */
			if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
				if ipPort != serverList[4] && ipPort !=serverList[5] {
					fmt.Println("Faulty Sender sends echo to ", ipPort )
					encoder.Encode(&data)
				}
			}
		} else {
			if sourceFault && source {
				if counter == 0 {
					encoder.Encode(&data)
				}
				if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
					if ipPort != serverList[4] && ipPort !=serverList[5] {
						fmt.Println("Faulty Sender sends echo to ", ipPort )
						encoder.Encode(&data)
					}
				}
			} else {
				fmt.Printf("Send Data to %+v\n",data)
				encoder.Encode(&data)
			}
		}
		counter++
	}
}



