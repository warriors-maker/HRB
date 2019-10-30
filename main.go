package main

import (
	"math/rand"
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
	//val1 := "abc"
	//val2 := "abcd"
	//
	//hashVal1 := HRBAlgorithm.ConvertBytesToString(HRBAlgorithm.Hash([]byte(val1)))
	//hashVal2 := HRBAlgorithm.ConvertBytesToString(HRBAlgorithm.Hash([]byte(val2)))
	//fmt.Println(hashVal1 == hashVal2)
	//fmt.Println(hashVal1)
	//fmt.Println(hashVal2)
}
