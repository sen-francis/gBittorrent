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

	reader := bufio.NewReader(file)
	// infinite loop
	count := 0
	for {

		b,err := reader.ReadByte()

		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println(err)
			break
		}

		if b == 'd' {
		// process the one byte b
			fmt.Println(b)
			fmt.Println(count)
		}
		
		if err != nil {
			// end of file
			break
		}
		count++
	}
	return ""
}
