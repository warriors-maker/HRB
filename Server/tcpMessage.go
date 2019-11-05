package Server

import "HRB/HRBAlgorithm"

type TcpMessage struct {
	Message HRBAlgorithm.Message
	ID      string
}
