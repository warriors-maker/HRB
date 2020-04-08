package Server

import (
	"HRB/HRBMessage"
)

type TcpMessage struct {
	Message HRBMessage.Message
	ID      string //which address I am sending to
}
