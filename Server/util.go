package Server

import (
	"bufio"
	"fmt"

	//"fmt"
	"log"
	"net"
	"os"
	//"strings"
)

//input: Path of the file
//output: list of server id, myId

func loadConfigServerFile(filePath string) []string{
	var addresses []string
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addr := scanner.Text()
		addresses = append(addresses, addr)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return addresses
}

func readTrustedFaulted(trustedPath, faultyPath string) ([]string, []string){
	trusted := loadConfigServerFile(trustedPath)
	faulty := loadConfigServerFile(faultyPath)
	fmt.Println(trusted, faulty)
	return trusted, faulty

}

func readServerListLocal(trustedPath, faultyPath string, index int) ([]string, string, bool){
	trusted, faulty := readTrustedFaulted(trustedPath, faultyPath)
	trustedCount = len(trusted)
	faultyCount = len(faulty)

	serverList := append(trusted, faulty...)
	myId := serverList[index];
	isFault := checkIsFaulty(myId, faulty)
	return serverList, myId, isFault
}


func readServerListNetwork(trustedPath, faultyPath string) ([]string, string, bool){
	trusted, faulty :=   readTrustedFaulted(trustedPath, faultyPath)
	serverList := append(trusted, faulty...)
	myId := getLocalIP()
	isFault := checkIsFaulty(myId, faulty)
	return serverList, myId, isFault
}


//Get my LocallIp
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func checkIsFaulty(myId string, faulty []string) bool{
	for _, v := range faulty {
		if v  == myId {
			return true
		}
	}
	return false
}