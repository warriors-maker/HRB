package ErasureCode

import (
	"fmt"
	"testing"
)

func TestDecodePermutation(t *testing.T) {

	dataShards := 3
	parityShards := 2
	val := "abcde"

	shards := Encode(val, dataShards, parityShards)
	fmt.Println(shards)
	result, _ := Decode(shards, dataShards, parityShards)
	fmt.Println("Result: " + result)

	vals := permutation(shards, dataShards, parityShards)

	count := 10
	results := []string {}
	for i := 0; i < count; i++ {
		results = append(results, val)
	}

	if len(vals) != len(results) {
		t.Errorf("Not same length %d, %d", len(vals), len(results))
	}

	for i := 0; i < len(vals); i++ {
		if vals[i] != results [i] {
			t.Errorf("Not same data %s, %s", vals[i], results[i])
		}
	}
}
