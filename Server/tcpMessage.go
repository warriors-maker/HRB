package Server

import "HRB/HRBAlgorithm"

type TcpMessage struct {
	Message HRBAlgorithm.Message
	ID      string //which address I am sending to
}
