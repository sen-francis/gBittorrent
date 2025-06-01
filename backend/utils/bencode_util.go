package utils

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func Decode(reader *bufio.Reader) (any, error) {
	b, err := reader.ReadByte()

	if err != nil {
		return "", err
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
	
	errorMessage := fmt.Sprintf("Could not decode bencoded text. Found unexpected char: %c", b)
	return "", errors.New(errorMessage)
}

func decodeDictionary(reader *bufio.Reader) (map[string]any, error) {
	dictionary := map[string]any{}
	b, err := reader.Peek(1)
	for err == nil && len(b) > 0 && b[0] != 'e' {
		key, err := decodeString(reader)
		if err != nil {
			return map[string]any{}, err
		}

		value, err := Decode(reader)
		if err != nil {
			return map[string]any{}, err
		}

		dictionary[key] = value
		b, err = reader.Peek(1)
		if err != nil {
			return map[string]any{}, err
		}
	}

	if b[0] == 'e' {
		reader.ReadByte()
	} else {
		return map[string]any{}, errors.New("Did not find dictionary terminator.")
	}

	return dictionary, nil
}

func decodeList(reader *bufio.Reader) ([]any, error) {
	list := []any{}
	b, err := reader.Peek(1)
	for err == nil && len(b) > 0 && b[0] != 'e' {
		val, err :=  Decode(reader)
		if err != nil {
			return []any{}, err
		}
		list = append(list, val)
		b, err = reader.Peek(1)
		if err != nil {
			return []any{}, err
		}
	}

	if b[0] == 'e' {
		reader.ReadByte()
	} else {
		return []any{}, errors.New("Did not find list terminator.")
	}


	return list, nil
}

func decodeString(reader *bufio.Reader) (string, error) {
	strLength, err := decodeInteger(reader, ':')
	if err != nil {
		return "", err
	}
	buf := make([]byte, strLength)
	n, err := reader.Read(buf)
	if (err != nil || n != int(strLength)) {
		return "", errors.New("Error while decoding byte string.")
	}
	return string(buf), nil
}

func decodeInteger(reader *bufio.Reader, delimeter byte) (int64, error) {
	numberStr, err := reader.ReadSlice(delimeter)
	
	if err != nil {
		return 0, errors.New("Did not find integer terminator")	
	}

	numberStr = numberStr[:len(numberStr) - 1]

	num, err := strconv.ParseInt(string(numberStr), 10, 64)

	if (err != nil) {
		return 0, errors.New("Could not convert string to integer")
	}

	return num, nil
}

