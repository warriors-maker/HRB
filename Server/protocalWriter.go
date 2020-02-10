package Server

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func TcpWriter(ch chan TcpMessage) {
	nets := strings.Split(MyId, ":")
	host := nets[0]
	externalPortNum, _ := strconv.Atoi(nets[1])

	benchmarkAddr := host + ":" + strconv.Itoa(externalPortNum + 500)
	fmt.Println("Protcal Channel for sending data to Benchmark " + benchmarkAddr)

	conn, err:= net.Dial("tcp", benchmarkAddr)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp", benchmarkAddr)
		time.Sleep(1*time.Second)
	}

	encoder := gob.NewEncoder(conn)
	for {
		data := <- ch
		//fmt.Printf("Protocal Send Data to Benchmark: %+v\n",data)
		encoder.Encode(&data)
	}
}