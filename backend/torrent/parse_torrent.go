package torrent

import (
	"bittorrent/backend/collections"
	"bittorrent/backend/utils"
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type FileInfo struct {
	Length int64
	Md5Sum string
	Path []string
}

type TorrentInfo struct {
	PieceLength int64
	Pieces []byte
	IsPrivate bool
	DirectoryName string
	FileInfoList []FileInfo
}

type TorrentMetainfo struct {
	Info TorrentInfo
	InfoHash [20]byte
	Announce string
	AnnounceList [][]string
	CreationDate int
	Comment string
	CreatedBy string
	Encoding string
	Size int64
}

func isFileListValid(fileList []map[string]any) bool {
	if len(fileList) == 0 {
		return false
	}

	for _, file := range fileList {
		if length, ok := file["length"]; ok {
			if _, ok := length.(int64); !ok {
				return false
			}
		} else {
			return false
		}
		if path, ok := file["path"]; ok {
			_, err := castFilePathList(path)
			if err != nil {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func castFilePathList(a any) ([]string, error) {
    v := reflect.ValueOf(a)
    switch v.Kind() {
    case reflect.Slice, reflect.Array:
        result := make([]string, v.Len())
        for i := range v.Len() {
			if val, ok := v.Index(i).Interface().(string); ok {
				result[i] = val
			} else {
				errorMsg := fmt.Sprintf("Unknown value in file path list: %s\n", reflect.ValueOf(val).Kind().String())
				return nil, errors.New(errorMsg) 
			}
		}
        return result, nil
    default:
		errorMsg := fmt.Sprintf("Provided value was not a slice: %s\n", v.Kind().String())
   		return nil, errors.New(errorMsg) 
	}
}

func isMultiFileInfoDictionaryValid(dictionary map[string]any) bool {
	if pieceLength, ok := dictionary["piece length"]; ok {
		if _, ok := pieceLength.(int64); !ok {
			return false
		}
	} else {
		return false
	}

	if pieces, ok := dictionary["pieces"]; ok {
		if _, ok := pieces.(string); !ok {
			return false
		}
	} else {
		return false
	}

	if name, ok := dictionary["name"]; ok {
		if _, ok := name.(string); !ok {
			return false
		}
	} else {
		return false
	}

	if files, ok := dictionary["files"]; ok {
		fileList, err := castFileList(files)
		if err != nil {
			return false	
		}
		if !isFileListValid(fileList) {
			return false	
		}
	} else {
		return false
	}

	return true
}


func castFileList(a any) ([]map[string]any, error) {
    v := reflect.ValueOf(a)
    switch v.Kind() {
    case reflect.Slice, reflect.Array:
        result := make([]map[string]any, v.Len())
        for i := range v.Len() {
			if val, ok := v.Index(i).Interface().(map[string]any); ok {
				result[i] = val
			} else {
				errorMsg := fmt.Sprintf("Unknown value in file list: %s\n", reflect.ValueOf(val).Kind().String())
				return nil, errors.New(errorMsg) 
			}
		}
        return result, nil
    default:
		errorMsg := fmt.Sprintf("Provided value was not a slice: %s\n", v.Kind().String())
   		return nil, errors.New(errorMsg) 
	}
}

func isSingleFileInfoDictionaryValid(dictionary map[string]any) bool {
	if pieceLength, ok := dictionary["piece length"]; ok {
		if _, ok := pieceLength.(int64); !ok {
			return false
		}
	} else {
		return false
	}

	if pieces, ok := dictionary["pieces"]; ok {
		if _, ok := pieces.(string); !ok {
			return false
		}
	} else {
		return false
	}

	if name, ok := dictionary["name"]; ok {
		if _, ok := name.(string); !ok {
			return false
		}
	} else {
		return false
	}

	if length, ok := dictionary["length"]; ok {
		if _, ok := length.(int64); !ok {
			return false
		}
	} else {
		return false
	}

	return true
}

func parseMultiFileInfoDictionary(dictionary map[string]any) (TorrentInfo, error) {
	var torrentInfo TorrentInfo
	if !isMultiFileInfoDictionaryValid(dictionary) {
		return torrentInfo, errors.New("Torrent Metainfo info dictionary formatted incorrectly")
	}

	torrentInfo.PieceLength = dictionary["piece length"].(int64)
	torrentInfo.Pieces = dictionary["pieces"].([]byte)
	torrentInfo.DirectoryName = dictionary["name"].(string)

	var fileInfoList []FileInfo
	files, _ := castFileList(dictionary["files"])
	for _, file := range files {
		filePath, _ := castFilePathList(file["path"])
		fileInfo := FileInfo { Length: file["length"].(int64), Path: filePath }

		if md5sum, ok := file["md5sum"]; ok {
			if md5sum, ok := md5sum.(string); ok {
				fileInfo.Md5Sum = md5sum
			}
		}
		
		fileInfoList = append(fileInfoList, fileInfo)
	}

	torrentInfo.FileInfoList = fileInfoList

	if isPrivate, ok := dictionary["private"]; ok {
		if isPrivate, ok := isPrivate.(int64); ok {
			torrentInfo.IsPrivate = isPrivate == 1
		}
	}

	return torrentInfo, nil
}

func parseSingleFileInfoDictionary(dictionary map[string]any) (TorrentInfo, error) {
	var torrentInfo TorrentInfo
	if !isSingleFileInfoDictionaryValid(dictionary) {
		return torrentInfo, errors.New("Torrent Metainfo info dictionary formatted incorrectly")
	}
	torrentInfo.PieceLength = dictionary["piece length"].(int64)
	torrentInfo.Pieces = dictionary["pieces"].([]byte)	
	fileInfo := FileInfo { Length: dictionary["length"].(int64), Path: []string{dictionary["name"].(string)} }	

	if md5sum, ok := dictionary["md5sum"]; ok {
		if md5sum, ok := md5sum.(string); ok {
			fileInfo.Md5Sum = md5sum
		}
	}

	torrentInfo.FileInfoList = []FileInfo {fileInfo}

	if isPrivate, ok := dictionary["private"]; ok {
		if isPrivate, ok := isPrivate.(int64); ok {
			torrentInfo.IsPrivate = isPrivate == 1
		}
	}

	return torrentInfo, nil
}

func parseInfoDictionary(dictionary map[string]any) (TorrentInfo, error) {
	_, isMultiFile := dictionary["files"]
	if isMultiFile {
		return parseMultiFileInfoDictionary(dictionary)
	} else {
		return parseSingleFileInfoDictionary(dictionary)
	}
}

func isTorrentMetainfoFileValid(dictionary map[string]any) bool {
	info, ok := dictionary["info"]
	if !ok {
		return false
	}

	if _, ok := info.(map[string]any); !ok {
		return false
	}

	announce, ok := dictionary["announce"]

	if !ok {
		return false
	}

	if _, ok := announce.(string); !ok {
		return false
	}

	return true
}

func splitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(substring)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func extractBencodedInfo(scanner *bufio.Scanner) (string, error) {
	scanner.Split(splitAt("4:info"))

	bencodedText := ""
	for scanner.Scan() {
		bencodedText = scanner.Text()
	}

	stack := collections.Stack[byte]{}
	bencodedInfo := ""
	bencodedTextArr := []byte(bencodedText)
	i := 0
	for i < len(bencodedTextArr) {
		if i != 0 && stack.IsEmpty() {
			break
		}

		c := bencodedTextArr[i]

		if c >= '0' && c <= '9' {
			numStr := ""
			for bencodedTextArr[i] != ':' {
				numStr += string(bencodedTextArr[i])
				i++
			}
			bencodedInfo += numStr + ":"
			i++
			num, err := strconv.ParseInt(numStr, 10, 0)
			if err != nil {
				return "", err
			}
			bencodedInfo += string(bencodedTextArr[i:i+int(num)])
			i += int(num)
			continue
		}

		switch c {
		case 'i':
			for bencodedTextArr[i] != 'e' {
				bencodedInfo += string(bencodedTextArr[i])
				i++
			}
			bencodedInfo += "e"
			i++
			continue
		case 'd', 'l':
			stack.Push(c)
		case 'e':
			stack.Pop()
		}

		bencodedInfo += string(c)
		i++
	}

	return bencodedInfo, nil
}


func (torrentInfo *TorrentInfo) getTotalTorrentBytes() int64 {
	totalBytes := int64(0)
	for _, fileInfo := range torrentInfo.FileInfoList {
		totalBytes += fileInfo.Length
	}

	return totalBytes
}

func ParseTorrentFile(filePath string) (TorrentMetainfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return TorrentMetainfo{}, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return TorrentMetainfo{}, err
	}

	reader := bufio.NewReaderSize(file, int(fileStat.Size()))

	result, err := utils.Decode(reader)
	if err != nil {
		return TorrentMetainfo{}, err
	}

	dictionary, ok := result.(map[string]any)
	if !ok || !isTorrentMetainfoFileValid(dictionary) {
		return TorrentMetainfo{}, errors.New("TorrentMetainfo file is invalid.")
	}

	infoDictionary := dictionary["info"].(map[string]any)

	info, err := parseInfoDictionary(infoDictionary)
	if err != nil {
		return TorrentMetainfo{}, err
	}

	var torrentMetainfo TorrentMetainfo
	torrentMetainfo.Info = info

	file.Seek(0, 0)
	bencodedInfo, err := extractBencodedInfo(bufio.NewScanner(file))
	if err != nil {
		return TorrentMetainfo{}, nil
	}
	torrentMetainfo.InfoHash = sha1.Sum([]byte(bencodedInfo))

	torrentMetainfo.Announce = dictionary["announce"].(string)

	if announceList, ok := dictionary["announce-list"]; ok {
		if announceList, ok := announceList.([][]string); ok {
			torrentMetainfo.AnnounceList = announceList
		}
	}

	if creationDate, ok := dictionary["creation date"]; ok {
		if creationDate, ok := creationDate.(int); ok {
			torrentMetainfo.CreationDate = creationDate
		}
	}

	if comment, ok := dictionary["comment"]; ok {
		if comment, ok := comment.(string); ok {
			torrentMetainfo.Comment = comment
		}
	}

	if createdBy, ok := dictionary["created by"]; ok {
		if createdBy, ok := createdBy.(string); ok {
			torrentMetainfo.CreatedBy = createdBy
		}
	}

	if encoding, ok := dictionary["encoding"]; ok {
		if encoding, ok := encoding.(string); ok {
			torrentMetainfo.Encoding = encoding
		}
	}
	
	torrentMetainfo.Size = torrentMetainfo.Info.getTotalTorrentBytes()

	return torrentMetainfo, nil
}
