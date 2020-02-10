package Server

//used in Benchmark
var internalReadChan chan TcpMessage
var internalWriteChan chan TcpMessage
var externalReadChan chan TcpMessage
var externalWriteChan map[string] messageChan

func BenchmarkStart() {
	initChannels()
	readSetup()
	writeSetup()
}

func initChannels() {
	internalReadChan = make (chan TcpMessage)
	internalWriteChan = make (chan TcpMessage)
	externalWriteChan = make (map[string] messageChan)
	externalReadChan = make(chan TcpMessage)
}

func readSetup() {
	go internalRead()
	go networkRead()
}

func writeSetup() {
	go internalWrite()
	go networkWrite()
}


func internalRead() {
	go internalReader(internalReadChan)
	for {
		data := <- internalReadChan
		sendTo := data.ID
		if sendTo == "" {
			for _, channel := range externalWriteChan {
				channel <- data
			}
		}  else {
			externalWriteChan[sendTo] <- data
		}
	}
}

func networkRead(){
	go ExternalTcpReader(externalReadChan, MyId)
	for {
		data := <- externalReadChan
		internalReadChan <- data
	}
}

func internalWrite() {
	internalWriter(protocalReadAddr, internalWriteChan)
}

func networkWrite() {
	//Responsible for writing to other servers
	for _, serverId := range serverList {
		externalWriteChan[serverId] = make(chan TcpMessage)
		go BenchmarkDeliver(serverId, externalWriteChan[serverId])
	}
}


func BenchmarkDeliver(ipPort string, ch chan TcpMessage) {
	ExternalTcpWriter(ipPort, ch)
}

