package main

import (
	"HRB/Server"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

/*
Run Program:
go run main.go [local/network] [algorithm] [Id[0,1...]] [f / ""]
go run main.go local 1
 */
func main() {
	//fmt.Println("Here are the four algorithms you can choose")
	//fmt.Println( "1. HashNonEquivocate")
	//fmt.Println( "2. HashEquivocate")
	//fmt.Println( "3. ECHashNonEquivocate")
	//fmt.Println( "4. ECHashEquivocate")

	var wg sync.WaitGroup
	wg.Add(1)

	argsWithoutProg := os.Args[1:]

	mode := argsWithoutProg[0]
	algorithm,_ := strconv.Atoi(argsWithoutProg[1])

	fmt.Println(mode)
	if mode == "network" {
		if len (argsWithoutProg) == 2 {
			soureceFault := false
			Server.Startup(-1, 5, soureceFault)
		} else {
			//SourceFault
			sourceFault := true
			Server.Startup(-1, 5, sourceFault)
		}
	} else if mode == "local" {
		// ./main local 5 idx
		if len (argsWithoutProg) == 3 {
			sourceFault := false
			idx, _ := strconv.Atoi(argsWithoutProg[2])
			Server.Startup(idx, algorithm, sourceFault)
		} else {
			//Source Fault
			sourceFault := true
			idx, _ := strconv.Atoi(argsWithoutProg[2])
			Server.Startup(idx, algorithm, sourceFault)
		}
	} else {
		fmt.Println("Invalid mode")
	}

	wg.Wait()

}
