package torrent

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"
	"encoding/binary"
)

const PSTR = "BitTorrent protocol"
const HANDSHAKE_LEN = 68
const RESERVED_LEN = 8

type Peer struct {
	PeerId string 
	IpAddress net.IP 
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
	
	err = peer.sendInterested()
	if err != nil {
		return err	
	}

	err = peer.waitForUnchoke()
	if err != nil {
		return err	
	}
	fmt.Printf("Connected to peer: %s\n", peer.String())
	return nil
}
const INTERESTED_LEN = 1
func (peer *Peer) generateInterested() []byte {
	interested := make([]byte, 5)
	binary.BigEndian.PutUint32(interested[0:4], uint32(INTERESTED_LEN))
	interested[4] = byte(INTERESTED)
	return interested
}

func (peer *Peer) sendInterested() error {
	interested := peer.generateInterested()
	_, err := peer.conn.Write(interested)
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}
	return nil	
}

func (peer *Peer) waitForUnchoke() error {
	for {
		deadline := time.Now().Add(2 * time.Minute)
		rawMessage, err := peer.readWithDeadline(deadline) 
		if err != nil {
			return err	
		}
		message := parseMessage(rawMessage)
		switch message.messageType {
		case KEEP_ALIVE:
			continue	
		case UNCHOKE:
			return nil
		default:
			errorMsg := fmt.Sprintf("Recieve unknown message from peer while waiting for unchoke:%s, %s\n", message.String(), peer.String())
			return errors.New(errorMsg)
		}
	}
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
	
	peerHandshake, err := peer.read()	
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}

	if !peer.validateHandshake(peerHandshake, infoHash) {
		failureReason := "Handshake from peer failed validation. peerId or infoHash did not match expected values."
		return peer.handleConnectionFailure(failureReason)
	}

	return nil
}

func (peer *Peer) read() ([]byte, error) {
	deadline := time.Now().Add(2 * time.Minute)
	return peer.readWithDeadline(deadline)
}

func (peer *Peer) readWithDeadline(deadline time.Time) ([]byte, error) {
	err := peer.conn.SetReadDeadline(deadline)
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
	bitfield, err := peer.read()
	if err != nil {
		return peer.handleConnectionFailure(err.Error())	
	}

	peer.bitfield = bitfield
	return nil
}

func (peer *Peer) closeConnection() {
	peer.conn.Close()
	peer.IsActive = false
}

func (peer *Peer) handleConnectionFailure(failureReason string) error {
	errorMessage := fmt.Sprintf("%s Peer: %s", failureReason, peer.String())
	peer.closeConnection()
	return errors.New(errorMessage)
}

func (peer *Peer) String() string {
	return net.JoinHostPort(peer.IpAddress.String(), strconv.Itoa(int(peer.Port)))
}

func (peer *Peer) validateHandshake(handshake []byte, infoHash [20]byte) bool {
	pStrLen := int(handshake[0])
	if pStrLen <= 0 {
		return false	
	}
	
	infoHashStartIndex := 1 + pStrLen + RESERVED_LEN 
	receivedInfoHash := handshake[infoHashStartIndex: len(handshake) - 20]
	if !bytes.Equal(receivedInfoHash, infoHash[:]){
		return false	
	}
	
	peerIdStartIndex := infoHashStartIndex + len(infoHash)
	receivedPeerId := string(handshake[peerIdStartIndex:])
	if peer.PeerId == "" {
		peer.PeerId = receivedPeerId	
	}

	if receivedPeerId != peer.PeerId {
		return false	
	}

	return true
}

func (peer *Peer) GetFirstAvailablePieceIndex(pieceMap map[int][]byte) (int, error) {
	for pieceIndex := range pieceMap {
		if peer.hasPiece(pieceIndex) {
			return pieceIndex, nil
		}
	}	
	return -1, peer.handleConnectionFailure("Peer has no required pieces.")
}

func (peer *Peer) hasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex >= len(peer.bitfield) {
		return false	
	}
	return ((peer.bitfield[byteIndex] >> (8 - offset)) & 1) != 0
}

func (peer *Peer) updateBitfield(index int) error {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex >= len(peer.bitfield) {
		errorMsg := fmt.Sprintf("Peer sent invalid have message: %s", peer.String())
		return errors.New(errorMsg)	
	}
	peer.bitfield[byteIndex] |= 1 << (8 - offset)
	return nil
}


