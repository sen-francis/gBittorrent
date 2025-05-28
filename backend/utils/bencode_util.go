package utils

import (
	"bufio"
	"strconv"
)

func decode(reader *bufio.Reader) any {
	b, err := reader.ReadByte()

	if err != nil {
		return err
	}

	if b == 'd' {
		return DecodeDictionary(reader)
	} else if b == 'l' {
		return decodeList(reader)
	} else if b == 'i' {
		return decodeInteger(reader, 'e')
	} else if _, err := strconv.Atoi(string(b)); err == nil {
		return decodeString(reader)
	}

	return "blah"
}

func DecodeDictionary(reader *bufio.Reader) map[string]any {
	dictionary := map[string]any{}	
	for b, err := reader.Peek(1); err == nil && len(b) > 0 && b[0] != 'e'; b, err = reader.Peek(1) {
		key := decodeString(reader)
		value := decode(reader)
		dictionary[key] = value
	}

	return dictionary
}

func decodeList(reader *bufio.Reader) []any {
	list := []any{}
	for b, err := reader.Peek(1); err == nil && len(b) > 0 && b[0] != 'e'; b, err = reader.Peek(1) {
		list = append(list, decode(reader))	
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

