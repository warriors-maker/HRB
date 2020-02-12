package Server

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"strings"
)

func internalReader(ch chan TcpMessage) {
	nets := strings.Split(internalReadAddr, ":")
	port := nets[1]
	ln, _ := net.Listen("tcp",":"+port)
	//fmt.Println("Listening internally at " + host + ":" + port)

	conn, err := ln.Accept()
	//fmt.Println("Get internal connection from " + conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}
	go internalHandleConnection(conn, ch)

}

func internalHandleConnection(conn net.Conn, ch chan TcpMessage) {
	defer conn.Close()
	/*
		Register the concrete Type
	*/
	dec := gob.NewDecoder(conn)

	data := &TcpMessage{}
	for {
		//Receive data
		if err := dec.Decode(data); err != nil {
			if errconn := conn.Close(); errconn != nil {
				os.Exit(1)
			}
		}

		//nets := strings.Split(internalReadAddr, ":")
		//host := nets[0]
		//port := nets[1]

		//fmt.Printf("Benchmark Receiving %+v from Protocal\n", *data)
		ch <- *data

	}
	//fmt.Println("Connection closed")
}


