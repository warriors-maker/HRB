package HRBAlgorithm

import (
	"bytes"
	"fmt"
	"github.com/klauspost/reedsolomon"
)


func Encode(message string, dataShards, parityShards int) [][]byte{
	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil
	}

	shards, err := enc.Split([]byte(message))
	if err != nil {
		return nil
	}

	err = enc.Encode(shards)
	if err != nil {
		return nil
	}

	fmt.Println(len(shards))
	return shards
}



func Decode(shard [][]byte, dataShards, parityShards int)  (string, error){

	shards := make([][]byte, len(shard))
	for i := range shard {
		shards[i] = make([]byte, len(shard[i]))
		copy(shards[i], shard[i])
	}

	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return "", err;
	}

	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Reconstructing data")
		err = enc.Reconstruct(shards)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			return "", err;
		}
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted.")
			return "", err;
		}
	}

	var buf bytes.Buffer
	size := len(shards[0])* dataShards
	err = enc.Join(&buf,shards, size)
	if err != nil {
		return "", err
	}

	value := buf.String()
	return value, nil
}

func DecodePermutation(shards [][]byte, dataShards, parityShards int, f func(string)) {
	aux := make([][]byte, dataShards + parityShards)
	decodePermutateHelper(shards, aux, dataShards, parityShards, 0, 0, f)
}

func decodePermutateHelper(shards, aux [][]byte,
	dataShards, parityShards, level, offset int,
	f func(string) ) {


	if level == dataShards{
		fmt.Println("Data",aux)
		val, err := Decode(aux, dataShards, parityShards)
		if err == nil {
			f(val)
		}
		return
	}

	for i := offset; i < dataShards + parityShards; i++ {
		if shards[i] != nil {
			aux[i] = shards[i]
			decodePermutateHelper(shards, aux, dataShards, parityShards, level + 1, i + 1, f)
			aux[i] = nil
		}
	}


}

func permutation(shards [][]byte, dataShards, parityShards int) []string{
	fmt.Println("Here are the shards"  ,shards)
	vals := []string{}
	DecodePermutation(shards, dataShards, parityShards, func (v string) {
		vals = append(vals, v)
	})
	return vals
}


func PermutateStr(input, aux []string, dataShards, parityShards, level, offset int) {
	if level == dataShards {
		fmt.Println(aux)
		return
	}
	for i:= offset; i < dataShards + parityShards; i++ {
		aux[i] = input[i]
		PermutateStr(input, aux, dataShards, parityShards, level + 1, i + 1)
		aux[i] = ""
	}
}



