package Server

import (
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
			if counter % 2 == 0 {

			} else {
				fmt.Printf("Sending : %+v\n", data.Message)
				time.Sleep(7*time.Second)
				encoder.Encode(&data)
			}
		} else {
			fmt.Printf("Sending : %+v\n", data.Message)
			time.Sleep(7*time.Second)
			encoder.Encode(&data)
		}

		counter++
	}
}



