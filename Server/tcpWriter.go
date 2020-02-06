package Server

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func TcpWriter(ch chan TcpMessage) {
	fmt.Println("Channel for sending data to " + MyId)

	conn, err:= net.Dial("tcp",MyId)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp",MyId)
		time.Sleep(3*time.Second)
	}
	encoder := gob.NewEncoder(conn)
	for {
		data := <- ch
		encoder.Encode(&data)
	}
}

//ipPort: the targer ipAddress to write to
//func TcpWriter(ipPort string, ch chan TcpMessage) {
//	counter := 0
//	fmt.Println("Channel for sending data to " + ipPort)
//
//	conn, err:= net.Dial("tcp",ipPort)
//
//	//keep dialing until the server comes up
//	for err != nil {
//		conn, err= net.Dial("tcp",ipPort)
//		time.Sleep(10*time.Second)
//	}
//
//	encoder := gob.NewEncoder(conn)
//	for {
//		data := <-ch
//		if isFault {
//			if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
//				fmt.Println("Set data to null")
//				correct := data.Message
//				// Create a Faulty Message
//
//				faulty := HRBAlgorithm.ECHOStruct{Id:correct.GetId(), Data: data.Message.GetData()+ "1asdadadwa", SenderId:correct.GetSenderId(),
//					HashData:correct.GetHashData(), Round:correct.GetRound(), Header:HRBAlgorithm.ECHO}
//				data = TcpMessage{Message:faulty, ID:MyId}
//
//				fmt.Println(data.Message.GetData())
//
//				encoder.Encode(&data)
//			} else {
//				encoder.Encode(&data)
//			}
//			//if counter % 2 == 0 {
//			//
//			//} else {
//			//	fmt.Printf("Sending : %+v to %v\n", data.Message, ipPort)
//			//	//time.Sleep(7*time.Second)
//			//	encoder.Encode(&data)
//			//}
//
//			/*
//			Test the Extreme Case for EC Coding
//			 */
//
//			//if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
//			//	if ipPort != serverList[4] && ipPort !=serverList[5] {
//			//		fmt.Println("Faulty Sender sends echo to ", ipPort )
//			//		encoder.Encode(&data)
//			//	}
//			//}
//		} else {
//			if sourceFault && source {
//				if counter == 0 {
//					encoder.Encode(&data)
//				}
//				if data.Message.GetHeaderType() == HRBAlgorithm.ECHO {
//					if ipPort != serverList[4] && ipPort !=serverList[5] {
//						fmt.Println("Faulty Sender sends echo to ", ipPort )
//						encoder.Encode(&data)
//					}
//				}
//			} else {
//				//time.Sleep(5*time.Second)
//				fmt.Printf("Hey, Send Data to %+v\n",data)
//				encoder.Encode(&data)
//			}
//		}
//		counter++
//	}
//}

