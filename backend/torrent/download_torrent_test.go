package torrent

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func generatePieces(pieceCount int) []byte {
	piecesLength := PIECE_HASH_LEN * pieceCount
	pieces := make([]byte, piecesLength)
	rand.Read(pieces)
	return pieces
}

func TestGeneratePieceMap(t *testing.T) {
	const minPieces = 50
	const maxPieces = 100
    nBig, err := rand.Int(rand.Reader, big.NewInt(maxPieces - minPieces)) 
    if err != nil {
        panic(err)
    }
    pieceCount := int(nBig.Int64() + minPieces)
	torrentInfo := TorrentInfo {
		Pieces: generatePieces(pieceCount),
	}
	torrentMetainfo := TorrentMetainfo{
		Info: torrentInfo,	
	}
	pieceMap := torrentMetainfo.generatePieceMap()
	if len(pieceMap) != pieceCount {
		t.Errorf("Expected pieceMap to be of length: %d, recieved: %d", pieceCount, len(pieceMap))
	}
}

func TestWritePieceToSingleFile(t *testing.T) { 
	pwd, _ := os.Getwd()
	outputPath := filepath.Join(pwd, "..", "..", "test_output")
	const fileName = "test-file.txt"
	filePath := filepath.Join(outputPath, fileName)
	initFileContents := []byte("piece0 piece1 piece2 piece3 piece4 piece5 ")
	err := os.WriteFile(filePath, initFileContents, 0644)
	if err != nil {
		t.Error(err)	
	}
	fileInfoList := []FileInfo {
		{ Length: 56, Path: []string {fileName}}, 
	}
	torrentInfo := TorrentInfo {
		PieceLength: 7,
		FileInfoList: fileInfoList,
	}
	torrentMetainfo := TorrentMetainfo {
		Info: torrentInfo,	
	}
	piece6 := []byte("piece6 ")
	err = torrentMetainfo.writePieceToFiles(6, piece6, outputPath)
	if err != nil {
		t.Error(err)	
	}

	piece7 := []byte("piece7 ")
	err = torrentMetainfo.writePieceToFiles(7, piece7, outputPath) 
	if err != nil {
		t.Error(err)	
	}
	expectedFileContents := append(initFileContents, piece6...)
	expectedFileContents = append(expectedFileContents, piece7...)
	fileContents, err := os.ReadFile(filePath)
	if !bytes.Equal(expectedFileContents, fileContents) {
		t.Errorf("Expected file contents to be: %s, recieved: %s", expectedFileContents, fileContents)
	}
	os.Remove(filePath)
}

func TestWritePieceToMultiFiles(t *testing.T) {
	pwd, _ := os.Getwd()
	outputPath := filepath.Join(pwd, "..", "..", "test_output")
	const fileName1 = "test-file_1.txt"
	const fileName2 = "test-file_2.txt"
	filePath1 := filepath.Join(outputPath, fileName1)
	filePath2 := filepath.Join(outputPath, fileName2)
	initFileContents := []byte("block0 block1 block2 block3 block4 block5 ")
	err := os.WriteFile(filePath1, initFileContents, 0644)
	if err != nil {
		t.Error(err)	
	}
	fileInfoList := []FileInfo {
		{ Length: 49, Path: []string {fileName1}}, 
		{Length: 14, Path: []string {fileName2}},
	}
	torrentInfo := TorrentInfo {
		PieceLength: 14,
		FileInfoList: fileInfoList,
	}
	torrentMetainfo := TorrentMetainfo {
		Info: torrentInfo,	
	}
	block6 := []byte("block6 ")
	block7 := []byte("block7 ")
	piece4 := append(block6, block7...)
	err = torrentMetainfo.writePieceToFiles(3, piece4, outputPath)
	if err != nil {
		t.Error(err)	
	}

	expectedFile1Contents := append(initFileContents, block6...)
	file1Contents, err := os.ReadFile(filePath1)
	if !bytes.Equal(expectedFile1Contents, file1Contents) {
		t.Errorf("Expected file 1 contents to be: %s, recieved: %s", expectedFile1Contents, file1Contents)
	}
	os.Remove(filePath1)


	expectedFile2Contents := block7 
	file2Contents, err := os.ReadFile(filePath2)
	if !bytes.Equal(expectedFile2Contents, file2Contents) {
		t.Errorf("Expected file 2 contents to be: %s, recieved: %s", expectedFile2Contents, file2Contents)
	}
	os.Remove(filePath2)
}

