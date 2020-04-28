package Server

type TcpMessage struct {
	Message interface{}
	ID      string //which address I am sending to
}
