package utils

import (
	"bufio"
	"fmt"
	"strconv"
)

func Decode(reader *bufio.Reader) any {
	b, err := reader.ReadByte()

	if err != nil {
		return err
	}

	if b == 'd' {
		return decodeDictionary(reader)
	} else if b == 'l' {
		return decodeList(reader)
	} else if b == 'i' {
		return decodeInteger(reader, 'e')
	} else if b >= '0' && b <= '9' {
		reader.UnreadByte()
		return decodeString(reader)
	}

	return "blah"
}

func decodeDictionary(reader *bufio.Reader) map[string]any {
	dictionary := map[string]any{}
	b, err := reader.Peek(1)
	for err == nil && len(b) > 0 && b[0] != 'e' {
		key := decodeString(reader)
		if key == "pieces" {
			fmt.Println("hi")
		}
		value := Decode(reader)
		dictionary[key] = value
		b, err = reader.Peek(1)
	}

	if b[0] == 'e' {
		reader.ReadByte()
	}

	return dictionary
}

func decodeList(reader *bufio.Reader) []any {
	list := []any{}
	b, err := reader.Peek(1)
	for err == nil && len(b) > 0 && b[0] != 'e' {
		list = append(list, Decode(reader))
		b, err = reader.Peek(1)
	}

	if b[0] == 'e' {
		reader.ReadByte()
	}

	return list
}

func decodeString(reader *bufio.Reader) string {
	strLength := decodeInteger(reader, ':')
	buf := make([]byte, strLength)
	n, err := reader.Read(buf)
	if (err != nil || n != int(strLength)) {
		// TODO SEN: blow up
		return ""
	}
	return string(buf)
}

func decodeInteger(reader *bufio.Reader, delimeter byte) int64 {
	numberStr, err := reader.ReadSlice(delimeter)
	
	if err != nil {
		return 0	
	}

	numberStr = numberStr[:len(numberStr) - 1]

	num, err := strconv.ParseInt(string(numberStr), 10, 64)

	if (err != nil) {
		// TODO SEN add error
		return 0
	}

	return num
}

