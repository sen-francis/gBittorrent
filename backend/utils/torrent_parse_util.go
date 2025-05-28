package utils

import (
	"os"
	"bufio"
	"errors"
)

type FileInfo struct {
	Length int64
	Md5Sum string
	Path []string
}

type TorrentInfo struct {
	PieceLength int64
	Pieces string
	IsPrivate bool
	DirectoryName string
	FileInfoList []FileInfo
}

type TorrentMetainfo struct {
	Info TorrentInfo
	InfoHash string
	Announce string
	AnnounceList [][]string
	CreationDate int
	Comment string
	CreatedBy string
	Encoding string
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
			if _, ok := path.([]string); !ok {
				return false
			}
		} else {
			return false
		}
	}

	return true
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

	if length, ok := dictionary["files"]; ok {
		if fileList, ok := length.([]map[string]any); !ok {
			return false
		} else {
			isFileListValid(fileList)	
		}
	} else {
		return false
	}

	return true
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
	torrentInfo.Pieces = dictionary["pieces"].(string)
	torrentInfo.DirectoryName = dictionary["name"].(string)

	var fileInfoList []FileInfo
	files := dictionary["files"].([]map[string]any)
	for _, file := range files {
		fileInfo := FileInfo { Length: file["length"].(int64), Path: file["path"].([]string) }	
		
		if md5sum, ok := file["md5sum"]; ok {
			if md5sum, ok := md5sum.(string); ok {
				fileInfo.Md5Sum = md5sum
			}
		}

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
	torrentInfo.Pieces = dictionary["pieces"].(string)	
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
		return parseSingleFileInfoDictionary(dictionary)
	} else {
		return parseMultiFileInfoDictionary(dictionary)
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

func ParseTorrentFile(filePath string) (TorrentMetainfo, error) {
	file, err := os.Open(filePath)
	var torrentMetainfo TorrentMetainfo
	if err != nil {
		return torrentMetainfo, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return torrentMetainfo, err
	}

	reader := bufio.NewReaderSize(file, int(fileStat.Size()))
	dictionary, ok := Decode(reader).(map[string]any)
	if !ok || !isTorrentMetainfoFileValid(dictionary) {
		return torrentMetainfo, errors.New("TorrentMetainfo file is invalid.")
	}

	infoDictionary := dictionary["info"].(map[string]any)

	info, err := parseInfoDictionary(infoDictionary)

	if err != nil {
		return torrentMetainfo, err
	}

	torrentMetainfo.Info = info

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


	return torrentMetainfo, nil
}
