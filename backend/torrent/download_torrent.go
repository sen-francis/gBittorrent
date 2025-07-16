package torrent

import (
	"bittorrent/backend/collections"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func (torrentMetainfo *TorrentMetainfo) StartDownload()  {
	peerQueue := collections.Queue[Peer]{}
	peerCh := make(chan []Peer)
	torrentStateCh := make(chan TorrentState)
	defer close(peerCh)
	defer close(torrentStateCh)
	torrentState := TorrentState{
		Event: "started", 
		Downloaded: 0, 
		Left: torrentMetainfo.Size, 
		Uploaded: 0, 
		PeerId: GeneratePeerId(),
	}
	torrentStateCh <- torrentState
	mutex := sync.Mutex{
	}
	go torrentMetainfo.fetchPeers(peerCh, torrentStateCh)
	go func() {
		for peerList := range peerCh {
			for _, peer := range peerList {
				go func() {
					err := peer.Connect(torrentMetainfo.InfoHash)
					if err != nil {
						fmt.Printf("Failed to connect to peer: %s", peer.String())
					}
					mutex.Lock()
					peerQueue.Push(peer)
					mutex.Unlock()
				}()
			}
		}
	}()
	
	pieceMap := torrentMetainfo.generatePieceMap()
	for len(pieceMap) > 0 {
		if peerQueue.IsEmpty() {
			fmt.Println("No peers available")
			continue	
		}

		peer, _ := peerQueue.Pop()
		if !peer.IsActive {
			fmt.Printf("Removed flaky peer from queue: %s", peer.String())
			continue	
		}
		pieceIndex, err := peer.GetFirstAvailablePieceIndex(pieceMap)
		if err != nil {
			fmt.Printf("Removing peer from queue: %s", err.Error())
			continue
		}

		pieceHash := pieceMap[pieceIndex]
		delete(pieceMap, pieceIndex)
		go func() {
			downloadedPiece, err := peer.downloadPiece(torrentMetainfo.Info.PieceLength, pieceIndex)
			mutex.Lock()
			if err != nil {
				fmt.Printf("Error downloading piece: %s\n", err)
				pieceMap[pieceIndex] = pieceHash
				mutex.Unlock()
			}
			if errors.Is(err, CHOKE_ERR) {
				fmt.Printf("Choked by peer: %s\n", peer.String())			
				err = peer.waitForUnchoke()	
				if err != nil {
					fmt.Printf("Peer did not unchoke within 10 minutes. Discarding peer: %s, %s\n", err.Error(), peer.String())	
					peer.closeConnection()
					return
				}
			}
			if err != nil {
				fmt.Printf("Discarding peer: %s, %s\n", peer.String())
				peer.closeConnection()
				return
			}
			mutex.Lock()	
			if !torrentMetainfo.verifyPiece(pieceHash, downloadedPiece) {	
				fmt.Printf("Could not verify downloaded piece from peer. Discarding peer: %s\n", peer.String())
				peer.closeConnection()
				mutex.Unlock()
				return
			}
			err = torrentMetainfo.writePieceToFiles(pieceIndex, downloadedPiece)
			if err != nil {	
				fmt.Printf("Failed to write to file: \n", err.Error())
			}
			peerQueue.Push(peer)
			mutex.Unlock()
		}()
	}
}

func (torrentMetainfo *TorrentMetainfo) writePieceToFiles(pieceIndex int64, piece []byte) error {
	fileStartIndex := int64(0)
	pieceStartIndex := torrentMetainfo.Info.PieceLength * pieceIndex
	for _, fileInfo := range torrentMetainfo.Info.FileInfoList {
		fileEndIndex := fileStartIndex + fileInfo.Length
		if pieceStartIndex + torrentMetainfo.Info.PieceLength < fileEndIndex {
			// piece in one file
			file, err := os.OpenFile(strings.Join(fileInfo.Path, "/"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) 
			if err != nil {
				return err	
			}
			_, err = file.WriteAt(piece, pieceStartIndex - fileStartIndex)
			return err
		} else if pieceStartIndex < fileEndIndex {
			// piece spanning multiple files	
			file, err := os.OpenFile(strings.Join(fileInfo.Path, "/"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) 
			if err != nil {
				return err	
			}
			_, err = file.WriteAt(piece[:fileEndIndex - pieceStartIndex], pieceStartIndex - fileStartIndex)
			if err != nil {
				return err	
			}
			pieceStartIndex = fileEndIndex
			piece = piece[fileEndIndex - pieceStartIndex :]
		}
		
		fileStartIndex += fileInfo.Length
	}
	return nil
}

func (torrentMetainfo *TorrentMetainfo) verifyPiece(expectedPieceHash []byte, piece []byte) bool {
	downloadedPieceHash := sha1.Sum(piece)
	return bytes.Equal(downloadedPieceHash[:], expectedPieceHash)
}

const PIECE_HASH_LEN = 20

func (torrentMetainfo *TorrentMetainfo) generatePieceMap() map[int][]byte {
	pieceMap := make(map[int][]byte)

	for startIndex:= 0; startIndex < len(torrentMetainfo.Info.Pieces); startIndex += PIECE_HASH_LEN { 
		index := startIndex / PIECE_HASH_LEN
		pieceMap[index] = torrentMetainfo.Info.Pieces[startIndex: startIndex + PIECE_HASH_LEN]
	}

	return pieceMap
}

type PieceState struct {
	downloadedBlocks int
	totalBlocks int
	blockOffset int
	piece []byte
}

const MAX_REQUESTS = 5
const BLOCK_SIZE = 16384

var CHOKE_ERR = errors.New("Choked by peer") 

func (peer *Peer) waitForUnchoke() error {
	deadline := time.Now().Add(10 * time.Minute)
	rawMessage, err := peer.readWithDeadline(deadline) 
	if err != nil {
		return err	
	}
	message := parseMessage(rawMessage)
	if message.messageType != UNCHOKE {
		errorMsg := fmt.Sprintf("Recieve unknown message from peer while waiting for unchoke:%s, %s\n", message, peer.String())
		return errors.New(errorMsg)
	}
	
	return nil
}

const REQUEST_RETRY_LIMIT = 3

func (peer *Peer) downloadPiece(pieceLength int64, pieceIndex int) ([]byte, error) {
	pieceState := PieceState{
		blockOffset: 0,	
		downloadedBlocks: 0,
		totalBlocks: int(pieceLength / BLOCK_SIZE),
		piece: make([]byte, pieceLength),
	}

	for range MAX_REQUESTS {
		err := peer.requestBlock(pieceIndex, pieceState.blockOffset)
		if err != nil {
			return nil, err
		}
	}

	for pieceState.downloadedBlocks < pieceState.totalBlocks {
		rawMessage, err := peer.read()
		if err != nil {
			return nil, err	
		}
		message := parseMessage(rawMessage)
		switch message.messageType {
		case KEEP_ALIVE: 
			continue
		case CHOKE:
			return nil, CHOKE_ERR
		case UNCHOKE:
			// should not happen here	
		case INTERESTED:
		// TODO SEEDING
		case NOT_INTERESTED: 
		// TODO SEEDING
		case HAVE:
			index := int(binary.BigEndian.Uint16(message.payload))
			peer.updateBitfield(index)
		case BITFIELD:
			peer.bitfield = message.payload
		case REQUEST:
		// TODO SEEDING
		case PIECE:
			index := int(binary.BigEndian.Uint16(message.payload[:4]))
			if index != pieceIndex {
				errorMsg := fmt.Sprintf("Piece index mismatch from peer: %s", peer.String())
				return nil, errors.New(errorMsg)	
			}
			offset := int(binary.BigEndian.Uint16(message.payload[4:8]))
			block := message.payload[8:]
			pieceState.downloadedBlocks++
			copy(pieceState.piece[offset: offset + BLOCK_SIZE], block)
			err := peer.requestBlock(pieceIndex, pieceState.blockOffset)
			if err != nil {
				return nil, err	
			}
		case CANCEL:
		// TODO SEEDING
		case PORT: 
		// TODO DHT
		}
	}
	return pieceState.piece, nil
}

const REQUEST_LEN = 13
func (peer *Peer) requestBlock(pieceIndex int, blockOffset int) error {
	request := make([]byte, 17)
	binary.BigEndian.PutUint32(request[0:4], uint32(REQUEST_LEN))
	request[4] = byte(REQUEST)
	binary.BigEndian.PutUint32(request[5:9], uint32(pieceIndex))
	binary.BigEndian.PutUint32(request[9:13], uint32(blockOffset))
	binary.BigEndian.PutUint32(request[9:13], uint32(BLOCK_SIZE))
	var err error
	for range REQUEST_RETRY_LIMIT {			
		_, err = peer.conn.Write(request)
		if err == nil {
			return nil

		}
	}
	return err
}
