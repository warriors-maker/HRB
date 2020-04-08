package Server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

/*
Read from other nodes
 */
func ExternalTcpReader(ch chan TcpMessage, listeningIp string) {
	nets := strings.Split(listeningIp, ":")
	//host := nets[0]
	port := nets[1]
	ln, _ := net.Listen("tcp",":"+port)
	fmt.Println("Benchmark Listening externally at addr: " + ":" + port)

	for {
		conn, err := ln.Accept()
		//fmt.Println("Get external connection from " + conn.RemoteAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		// Each HandleConnection handle one connection with one node
		go ExternalHandleConnection(conn, ch)
	}
}

func ExternalHandleConnection(conn net.Conn, ch chan TcpMessage) {
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

		//fmt.Printf("Receiving from others externally %+v \n",data.Message)

		ch <- *data

	}
}

