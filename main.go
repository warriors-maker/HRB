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
	if len(argsWithoutProg) == 0 {
		Server.NetworkModeStartup()
	} else {
		fmt.Println("Local Mode")
		idx, _ := strconv.Atoi(argsWithoutProg[0])
		Server.LocalModeStartup(idx)
	}
	//The server should never sleep
	wg.Wait()
}
