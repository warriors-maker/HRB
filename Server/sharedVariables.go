package Server

import (
	"strconv"
	"strings"
)

var default_benchmark_port = "9000"

var serverList []string
var MyId string //basically the IP of individual server (Note this is Benchmark port )
var isFault bool //check whether I should be the faulty based on configuration
var faultyCount int
var trustedCount int
var isLocalMode bool //indicate whether this is a local mode
var source bool //A flag to indicate whether I am the sender
var algorithm int
var isSourceFault bool
var faultyList []string
var trustedList []string

type messageChan chan TcpMessage

var internalReadAddr string
var protocalReadAddr string

func InitSharedVariables(index int) {
	//network mode = -1, local mode >= 0
	yamlStruct := decodeYamlFile()

	/*
		Initialization of the Global Variable
	*/

	faultyList = yamlStruct.Faulty
	trustedList = yamlStruct.Trusted
	serverList = append(trustedList, faultyList...)

	trustedCount = len(trustedList)
	faultyCount = len(faultyList)

	//-1: Network Benchmark
	// >0: Local Benchmark:
	// where index = index of the serverList => MyId
	// index 0 is always the source
	if index == -1 {
		isLocalMode = false
		myHostAddr := getLocalIP()
		MyId = myHostAddr+":" + default_benchmark_port
	} else {
		isLocalMode = true
		if index == 0 {
			source = true
		} else {
			source = false
		}
		MyId = serverList[index]
	}
	checkIsFault()
	algorithm = yamlStruct.Algorithm
	isSourceFault = yamlStruct.Source_Byzantine
	initOtherAddr()
}

func checkIsFault() {
	for _, id := range faultyList {
		if MyId == id {
			isFault = true
		}
	}
}

func initOtherAddr() {
	s := strings.Split(MyId, ":")

	hostAddr := s[0]
	externalPort, _ := strconv.Atoi(s[1])

	internalReadPort := strconv.Itoa(externalPort + 500)
	protocalReadPort := strconv.Itoa(externalPort + 1000)

	internalReadAddr = hostAddr + ":" + internalReadPort
	protocalReadAddr = hostAddr + ":" + protocalReadPort
}