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

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	argsWithoutProg := os.Args[1:]

	mode := argsWithoutProg[0]
	fmt.Println(mode)
	if mode == "network" {
		if len (argsWithoutProg) == 1 {

		} else {
			//SourceFault
		}
	} else if mode == "local" {
		if len (argsWithoutProg) == 2 {
			sourceFault := false
			idx, _ := strconv.Atoi(argsWithoutProg[1])
			Server.LocalModeStartup(idx, sourceFault)
		} else {
			//Source Fault
			sourceFault := true
			idx, _ := strconv.Atoi(argsWithoutProg[1])
			Server.LocalModeStartup(idx, sourceFault)
		}
	} else {
		fmt.Println("Invalid mode")
	}

	wg.Wait()


}
