package Network

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func TcpReader(ch chan message) {
	ln, _ := net.Listen("tcp",":"+MyPort)

	for {
		conn, err := ln.Accept()
		fmt.Println("Get connection from " + conn.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		// Each HandleConnection handle one connection with one node
		go handleConnection(conn, ch)
	}
}

func handleConnection(conn net.Conn, ch chan message) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	data := &message{}
	for {
		//Receive data
		dec.Decode(data)

		ch <- *data

		//Send data
		time.Sleep(3*time.Second)

		//ack := message{}
		//ack.Flag = 5
		//err := encoder.Encode(&data)
		//fmt.Println("Prepare sending flag " + strconv.Itoa(data.Flag) + " from " + strconv.Itoa(data.ID) + " detected ",data.Detected)
		//checkErr(err)

	}
	fmt.Println("Connection closed")
}