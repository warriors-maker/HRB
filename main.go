package main

import (
	"HRB/Server"
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
	var wg sync.WaitGroup
	wg.Add(1)

	argsWithoutProg := os.Args[1:]

	// If the running argument has an extra index, then it is local mode
	// where index represents what port local machine is using

	if len(argsWithoutProg) == 0 {
		Server.Startup(-1)
	} else {
		index,_ := strconv.Atoi(argsWithoutProg[0])
		Server.Startup(index)
	}
	wg.Wait()

}
