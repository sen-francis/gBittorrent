package utils

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"errors"
)

/*func check(e error) {
    if e != nil {
        panic(e)
    }
}
*/

func ParseTorrentFile(filePath string) string {
	// TODO SEN: at some point we should move away from reading the file into memory
	file, err := os.Open(filePath)

	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println(err)
		return ""
	}

	fileStat, err := file.Stat()
	if err != nil {
	}

	reader := bufio.NewReaderSize(file, int(fileStat.Size()))

	b, err := reader.Peek(15)
	fmt.Println(string(b))

	// infinite loop
	tmp := Decode(reader)
	fmt.Println(tmp)

	return ""
}
