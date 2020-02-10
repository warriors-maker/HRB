package HRBAlgorithm

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"math/rand"
)

/*
Hash Function
 */

const SecretKey string = "secret"

func Hash(input []byte) []byte{
	hmac512 := hmac.New(sha512.New, []byte(SecretKey))
	hmac512.Write(input)
	hash := hmac512.Sum(nil)
	return hash
}

//First is the expected data bytes
// Second is sent from others
func ValidData(data, hashData []byte) bool {
	expectedHash := Hash(data)
	return hmac.Equal(expectedHash, hashData)
}

func ValidHash(expectedHash, hashData []byte) bool{
	return hmac.Equal(expectedHash, hashData)
}

//Since key does not support byte[]
func ConvertBytesToString(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func ConvertStringToBytes(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}