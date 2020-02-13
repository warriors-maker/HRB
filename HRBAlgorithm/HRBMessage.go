package HRBAlgorithm


type TcpHeader int

const (
	MSG  TcpHeader = 0
	ECHO TcpHeader = 1
	ACC  TcpHeader = 2
	REQ  TcpHeader = 3
	FWD  TcpHeader = 4
	BIN  TcpHeader = 5
	Stat  TcpHeader = 6
	MSG_OPT TcpHeader = 7
	ECHO_OPT TcpHeader = 8
	ACC_OPT TcpHeader = 9
	REQ_OPT TcpHeader = 10
	FWD_OPT TcpHeader = 11
)

type Message interface {
	GetHeaderType() TcpHeader
	//GetHashData() []byte
	//GetData() []byte
	//GetId() string
	//GetRound() int
	//GetSenderId() string
	//SetSenderIdNull()
	//SetSenderId(string)
}

type FWDStruct struct {
	Header TcpHeader
	Data   []byte
	HashData []byte
	Round int
	Id string
	SenderId string
}

type MSGStruct struct {
	Header TcpHeader
	Data   []byte
	HashData []byte
	Round int
	Id string
	SenderId string
}

type ECHOStruct struct {
	Header TcpHeader
	HashData []byte
	Data []byte
	Round int
	Id string
	SenderId string
}

type ACCStruct struct {
	Header TcpHeader
	HashData []byte
	Round int
	Id string
	SenderId string
}

type REQStruct struct {
	Header TcpHeader
	HashData []byte
	Round int
	Id string
	SenderId string
}

type PrepareSend struct {
	M Message
	SendTo string
	Stat Stats
}

type Binary struct {
	Header TcpHeader
	HashData []byte
	Round int
	SenderId string
	Id string
}

type StatStruct struct {
	Header TcpHeader
	Round int
	Id string
}

func (s FWDStruct) GetHeaderType() TcpHeader{
	return s.Header
}

func (s MSGStruct) GetHeaderType() TcpHeader{
	return s.Header
}

func (s REQStruct) GetHeaderType() TcpHeader{
	return s.Header
}

func (s ACCStruct) GetHeaderType() TcpHeader{
	return s.Header
}

func (s ECHOStruct) GetHeaderType() TcpHeader{
	return s.Header
}

func (s Binary) GetHeaderType() TcpHeader{
	return s.Header
}

func (s StatStruct) GetHeaderType() TcpHeader{
	return s.Header
}


