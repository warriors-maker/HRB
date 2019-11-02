package Network

import (
	"ByzantineConsensusAlgorithm/Utility"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

//ipPort: the targer ipAddress to write to
func TcpWriter(ipPort string, ch chan TcpMessage) {
	fmt.Println("Channel for sending data to " + ipPort)

	conn, err:= net.Dial("tcp",ipPort)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp",ipPort)
		time.Sleep(10*time.Second)
	}

	encoder := gob.NewEncoder(conn)
	//dec := gob.NewDecoder(conn)
	for {
		data := <-ch
		fmt.Printf("Sending : %+v\n", data.Message)

		err := encoder.Encode(&data)
		Utility.CheckErr(err)
	}
}



