package HRBAlgorithm


type TcpHeader int

const (
	MSG  TcpHeader = 0
	ECHO TcpHeader = 1
	ACC  TcpHeader = 2
	REQ  TcpHeader = 3
	FWD  TcpHeader = 4

)

type Message interface {
	GetHeaderType() TcpHeader
	GetHashData() string
	GetData() string
	GetId() string
	GetRound() int
	GetSenderId() string
	SetSenderIdNull()
	SetSenderId(string)
}

type FWDStruct struct {
	Header TcpHeader
	Data   string
	Round int
	Id string
	SenderId string
}


type MSGStruct struct {
	Header TcpHeader
	Data   string
	Round int
	Id string
	SenderId string
}

type ECHOStruct struct {
	Header TcpHeader
	HashData string
	Round int
	Id string
	SenderId string
}

type ACCStruct struct {
	Header TcpHeader
	HashData string
	Round int
	Id string
	SenderId string
}

type REQStruct struct {
	Header TcpHeader
	HashData string
	Round int
	Id string
	SenderId string
}


type ReqChanStruct struct {
	M Message
	SendTo string
}

/*
Implement the interface
 */

func (m FWDStruct) GetHeaderType() TcpHeader{
	return m.Header
}

func (m MSGStruct) GetHeaderType() TcpHeader{
	return m.Header
}

func (m ECHOStruct) GetHeaderType() TcpHeader{
	return m.Header
}

func (m ACCStruct) GetHeaderType() TcpHeader{
	return m.Header
}

func (m REQStruct) GetHeaderType() TcpHeader{
	return m.Header
}


/*
Return the Value contained in the message for real message
 */

func (m FWDStruct) GetData() string{
	return m.Data
}

func (m MSGStruct) GetData() string {
	return m.Data
}

func (m ECHOStruct) GetData() string {
	return ""
}

func (m ACCStruct) GetData() string {
	return ""
}

func (m REQStruct) GetData() string {
	return ""
}

/*
Return the Value contained in the message for real message
*/

func (m FWDStruct) GetHashData() string{
	return ""
}

func (m MSGStruct) GetHashData() string{
	return ""
}

func (m ECHOStruct) GetHashData() string{
	return m.HashData
}

func (m ACCStruct) GetHashData() string{
	return m.HashData
}

func (m REQStruct) GetHashData() string{
	return m.HashData
}


/*
Get the senderID of this message
*/

func (m FWDStruct) GetId() string{
	return m.Id
}

func (m MSGStruct) GetId() string{
	return m.Id
}

func (m ECHOStruct) GetId() string{
	return m.Id
}

func (m ACCStruct) GetId() string{
	return m.Id
}

func (m REQStruct) GetId() string{
	return m.Id
}


/*
Get the Round
*/

func (m FWDStruct) GetRound() int{
	return m.Round
}

func (m MSGStruct) GetRound() int{
	return m.Round
}

func (m ECHOStruct) GetRound() int{
	return m.Round
}

func (m ACCStruct) GetRound() int{
	return m.Round
}

func (m REQStruct) GetRound() int{
	return m.Round
}



func (m FWDStruct) GetSenderId() string{
	return m.SenderId
}

func (m MSGStruct) GetSenderId() string{
	return m.SenderId
}

func (m ECHOStruct) GetSenderId() string{
	return m.SenderId
}

func (m ACCStruct) GetSenderId() string{
	return m.SenderId
}

func (m REQStruct) GetSenderId() string{
	return m.SenderId
}


func (m FWDStruct) SetSenderIdNull() {
	m.SenderId = ""
}

func (m MSGStruct) SetSenderIdNull() {
	m.SenderId = ""
}

func (m ECHOStruct) SetSenderIdNull() {
	m.SenderId = ""
}

func (m ACCStruct) SetSenderIdNull() {
	m.SenderId = ""
}

func (m REQStruct) SetSenderIdNull() {
	m.SenderId = ""
}




func (m FWDStruct) SetSenderId(id string) {
	m.SenderId = id
}

func (m MSGStruct) SetSenderId(id string) {
	m.SenderId = id
}

func (m ECHOStruct) SetSenderId(id string) {
	m.SenderId = id
}

func (m ACCStruct) SetSenderId(id string) {
	m.SenderId = id
}

func (m REQStruct) SetSenderId(id string) {
	m.SenderId = id
}
