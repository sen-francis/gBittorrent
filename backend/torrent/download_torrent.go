package torrent

import (
	"bittorrent/backend/collections"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
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

		piece := pieceMap[pieceIndex]
		delete(pieceMap, pieceIndex)
		go func() {
			err := peer.downloadPiece(torrentMetainfo.Info.PieceLength, pieceIndex)
			mutex.Lock()
			if err != nil {
				pieceMap[pieceIndex] = piece
			}
			if errors.Is(err, CHOKE_ERR) {
				fmt.Printf("Choked by peer: %s\n", peer.String())			
				// todo sen wait for unchoke
			}
			// todo sen: verify piece with piece hash
			peerQueue.Push(peer)
			mutex.Unlock()
		}()
	}
}


func (torrentMetainfo *TorrentMetainfo) generatePieceMap() map[int][]byte {
	pieceMap := make(map[int][]byte)

	for startIndex:= 0; startIndex < len(torrentMetainfo.Info.Pieces); startIndex += PIECE_HASH_LEN { 
		// todo sen: this might break on last index
		index := startIndex / PIECE_HASH_LEN
		pieceMap[index] = torrentMetainfo.Info.Pieces[startIndex: startIndex + PIECE_HASH_LEN]
	}

	return pieceMap
}

type PieceState struct {
	downloadedBlocks []bool
	blockOffset int
	currentRequests int
}

const MAX_REQUESTS = 5
const BLOCK_SIZE = 16384

var CHOKE_ERR = errors.New("Choked by peer") 

func (peer *Peer) waitForUnchoke() error {
	rawMessage, err := peer.readWithDeadline() // todo wait longer than 2 mins for unchoke
	// can we recieve anything other than unchoke while choked?
}

func (peer *Peer) downloadPiece(pieceLength int64, pieceIndex int) error {
	pieceState := PieceState{
		blockOffset: 0,	
		currentRequests: 0,
		downloadedBlocks: make([]bool, pieceLength / BLOCK_SIZE),
	}

	for pieceState.currentRequests < MAX_REQUESTS {
		// send request
	}
	for pieceState.currentRequests > 0 {
		rawMessage, err := peer.readWithDeadline()	
		if err != nil {
			return err	
		}
		message := parseMessage(rawMessage)
		switch message.messageType {
		case KEEP_ALIVE: 
			continue
		case CHOKE:
			return CHOKE_ERR
		case UNCHOKE:
			return nil	
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
			offset := int(binary.BigEndian.Uint16(message.payload[4:8]))
			block := message.payload[8:]
			return PIECE, nil
		case CANCEL:
		// TODO SEEDING
		case PORT: 
		// TODO DHT
		}
	}
	return nil
}
