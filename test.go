package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	str := "my_secret_sign_key"

	// Encode the string to Base64
	encodedString := base64.StdEncoding.EncodeToString([]byte(str))
	fmt.Println("Encoded:", encodedString)

	// Decode the Base64 string back to the original
	decodedString, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	fmt.Println("Decoded:", string(decodedString))
}
