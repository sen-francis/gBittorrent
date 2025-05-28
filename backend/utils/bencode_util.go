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
		return decodeInteger(reader)
	} else if _, err := strconv.Atoi(string(b)); err == nil {
		return decodeString(reader)
	}

	return "blah"
}

func DecodeDictionary(reader *bufio.Reader) map[string]any {
		
	return map[string]any{}
}

func decodeList(reader *bufio.Reader) []any {
	return []any{}
}

func decodeString(reader *bufio.Reader) string {
	return ""
}

func decodeInteger(reader *bufio.Reader) int {
	numberStr := ""
	reader.ReadString('e')
	for b, err := reader.ReadByte(); b != 'e' && err == nil; b, err = reader.ReadByte() {
		numberStr += string(b)	
	}

	num, err := strconv.Atoi(string(numberStr))

	if (err != nil) {
		// TODO SEN add error
		return 0
	}

	return num
}

