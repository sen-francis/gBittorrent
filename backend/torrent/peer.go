package torrent

import (
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"time"
	"errors"
)

const PSTR_LEN = 19
const PSTR = "BitTorrent protocol"
const HANDSHAKE_LEN = 68

type Peer struct {
	PeerId string 
	IpAddress string
	Port uint 
}

func GeneratePeerId() string {
	random := ""
	for range 12 {
		random += strconv.Itoa(rand.IntN(10))
	}
	return AZ_CLIENT_PREFIX + random
}

func (peer *Peer) Connect(infoHash string) (error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return err
	}

	err = peer.handshake(conn, infoHash)
	if err != nil {
		return err	
	}
	return nil
}

func (peer *Peer) generateHandshake(infoHash string) []byte {
	handshake := make([]byte, HANDSHAKE_LEN)
	handshake[0] = byte(PSTR_LEN)
	index := 1	
	index += copy(handshake[index:], PSTR)
	index += copy(handshake[index:], make([]byte, 8))
	index += copy(handshake[index:], []byte(infoHash))
	copy(handshake[index:], []byte(GeneratePeerId()))
	return handshake
}

func (peer *Peer) handshake(conn net.Conn, infoHash string) error {
	handshake := peer.generateHandshake(infoHash)
	_, err := conn.Write(handshake)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to connect to peer at %s", peer.String())
		conn.Close()
		return errors.New(errorMessage)
	}
	
	start := time.Now()	
	reply := make([]byte, 1024)

	for time.Since(start) < 2 * time.Minute {
		// todo sen: wait up to a minute for an unchoke response
		_, err = conn.Read(reply)
		if err != nil {
			conn.Close()
			println("Write to server failed:", err.Error())
		}
	}

}

func (peer *Peer) String() string {
	return peer.IpAddress + ":" + strconv.FormatUint(uint64(peer.Port), 10)
}
