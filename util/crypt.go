package util

import "encoding/hex"

func ToHex(s string) ([]byte, error) {
	hexByte := []byte("")
	//this if is for make the string len even
	if len(s)%2 == 1 {
		s = "0" + s
	}
	hexByte, err := hex.DecodeString(s)
	if err != nil {
		return hexByte, err
	}
	return hexByte, nil
}

func ReverseByte(input []byte) []byte {
	len := len(input)
	var output []byte
	for i := len - 1; i >= 0; i-- {
		output = append(output, input[i])
	}
	return output
}
