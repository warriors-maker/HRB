package Network

import (
	"ByzantineConsensusAlgorithm/Utility"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

//ipPort: the targer ipAddress to write to
func TcpWriter(ipPort string, ch chan message) {
	fmt.Println("Channel for sending data to " + ipPort)

	conn, err:= net.Dial("tcp",ipPort)
	for err != nil {
		conn, err= net.Dial("tcp",ipPort)
		time.Sleep(10*time.Second)
	}

	encoder := gob.NewEncoder(conn)
	//dec := gob.NewDecoder(conn)
	for {
		data := <-ch
		data.ID = myID
		err := encoder.Encode(&data)
		Utility.CheckErr(err)
	}
}

