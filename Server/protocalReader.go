package Server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

/*
Protocol reading from benchmark
 */

func TcpReader(ch chan TcpMessage, listeningIp string) {
	portNum, _ := strconv.Atoi(strings.Split(listeningIp, ":")[1])
	port := strconv.Itoa(portNum + 1000)
	ln, _ := net.Listen("tcp",":"+port)

	conn, err := ln.Accept()
	fmt.Println("Protocal Listening internally at ",port)
	if err != nil {
		log.Fatal(err)
	}
	// Each HandleConnection handle one connection with one node
	go handleConnection(conn, ch)
}

func handleConnection(conn net.Conn, ch chan TcpMessage) {
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

		fmt.Printf("Protocal Receiving from Benchmark %+v\n",data.Message)

		ch <- *data

	}
	//fmt.Println("Connection closed")
}