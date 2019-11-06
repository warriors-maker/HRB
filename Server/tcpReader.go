package Server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func TcpReader(ch chan TcpMessage, listeningIp string) {
	port := strings.Split(listeningIp, ":")[1]
	ln, _ := net.Listen("tcp",":"+port)

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

		fmt.Printf("Receiving %+v\n",data.Message)

		ch <- *data

		//Send data
		//time.Sleep(3*time.Second)

		//ack := message{}
		//ack.Flag = 5
		//err := encoder.Encode(&data)
		//fmt.Println("Prepare sending flag " + strconv.Itoa(data.Flag) + " from " + strconv.Itoa(data.ID) + " detected ",data.Detected)
		//checkErr(err)

	}
	fmt.Println("Connection closed")
}