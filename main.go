package main

import (
	"HRB/Server"
	"fmt"
	"os"
	"strconv"
	"sync"
)


/*
Run Program:
local Mode: go run main.go index
cluster Mode: go run main.go
 */

/*
Local mode:

benchmark reading from other nodes port: 	inputPort
benchmark reading from protocol port:		inputPort+500
protocol reading from benchmark port:		inputPort+1000


 */

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	argsWithoutProg := os.Args[1:]

	// If the running argument has an extra index, then it is local mode
	// where index represents what port local machine is using

	if len(argsWithoutProg) == 0 {
		Server.InitSharedVariables(-1) //cluster mode
	} else {
		fmt.Println(argsWithoutProg[0])
		index,_ := strconv.Atoi(argsWithoutProg[0]) //locolhost mode
		Server.InitSharedVariables(index)
	}

	//start the benchmark
	Server.BenchmarkStart()
	Server.ProtocalStart()

	wg.Wait()
}

