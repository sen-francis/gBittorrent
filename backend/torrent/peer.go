package torrent

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"
)

const PSTR = "BitTorrent protocol"
const HANDSHAKE_LEN = 68
const RESERVED_LEN = 8

type Peer struct {
	PeerId string 
	IpAddress string
	Port uint
	conn net.Conn
	bitfield []byte
	IsActive bool
}

func GeneratePeerId() string {
	random := ""
	for range 12 {
		random += strconv.Itoa(rand.IntN(10))
	}
	return AZ_CLIENT_PREFIX + random
}

func (peer *Peer) Connect(infoHash [20]byte) (error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return err
	}
	peer.conn = conn

	err = peer.handshake(infoHash)
	if err != nil {
		return err	
	}

	err = peer.receiveBitfield()
	if err != nil {
		return err
	}

	return nil
}

func (peer *Peer) generateHandshake(infoHash [20]byte) []byte {
	handshake := make([]byte, HANDSHAKE_LEN)
	handshake[0] = byte(len(PSTR))
	index := 1	
	index += copy(handshake[index:], PSTR)
	index += copy(handshake[index:], make([]byte, RESERVED_LEN))
	index += copy(handshake[index:], infoHash[:])
	copy(handshake[index:], []byte(GeneratePeerId()))
	return handshake
}

func (peer *Peer) handshake(infoHash [20]byte) error {
	handshake := peer.generateHandshake(infoHash)
	_, err := peer.conn.Write(handshake)
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}
	
	peerHandshake, err := peer.readWithDeadline()	
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}

	if !peer.validateHandshake(peerHandshake, infoHash) {
		failureReason := "Handshake from peer failed validation. peerId or infoHash did not match expected values."
		return peer.handleConnectionFailure(failureReason)
	}

	return nil
}

func (peer *Peer) readWithDeadline() ([]byte, error) {
	waitTime := time.Now().Add(2 * time.Minute)
	err := peer.conn.SetReadDeadline(waitTime)
	if err != nil {
		return nil, err 
	}

	buf := make([]byte, 1024)
	n, err := peer.conn.Read(buf)
	if err != nil || n <= 0 {
		return nil, err
	}

	return buf[:n], nil
}

func (peer *Peer) receiveBitfield() (error) {
	bitfield, err := peer.readWithDeadline()
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}

	peer.bitfield = bitfield
	return nil
}


func (peer *Peer) handleConnectionFailure(failureReason string) error {
	errorMessage := fmt.Sprintf("%s. Peer: %s", failureReason, peer.String())
	peer.conn.Close()
	peer.IsActive = false
	return errors.New(errorMessage)
}

func (peer *Peer) String() string {
	return peer.IpAddress + ":" + strconv.FormatUint(uint64(peer.Port), 10)
}

func (peer *Peer) validateHandshake(handshake []byte, infoHash [20]byte) bool {
	pStrLen := int(handshake[0])
	if pStrLen <= 0 {
		return false	
	}
	
	infoHashStartIndex := 1 + pStrLen + RESERVED_LEN 
	receivedInfoHash := handshake[infoHashStartIndex: len(handshake) - 20]
	if bytes.Equal(receivedInfoHash, infoHash[:]){
		return false	
	}
	
	peerIdStartIndex := infoHashStartIndex + len(infoHash)
	receivedPeerId := string(handshake[peerIdStartIndex:])
	if receivedPeerId != peer.PeerId {
		return false	
	}

	return true
}

func (peer *Peer) GetFirstAvailablePiece(pieceList *list.List) []byte {
	for e := pieceList.Front(); e != nil; e = e.Next() {
		piece, _ := e.Value.([]byte)	
		if peer.hasPiece(piece) {
			return piece
		}
	}
	return make([]byte, 1)
}

func (peer *Peer) hasPiece(piece []byte) bool {
	return false
}
