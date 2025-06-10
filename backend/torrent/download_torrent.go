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

func (peer *Peer) handshake(conn net.Conn, infoHash string) error {
	[]byte
	_, err := conn.Write([]byte(strEcho))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to connect to peer at %s", peer.String())
		conn.Close()
		return errors.New(errorMessage)
	}

	reply := make([]byte, 1024)
	
	// todo sen: wait up to a minute for an unchoke response
	_, err = conn.Read(reply)
	if err != nil {
		conn.Close()
		println("Write to server failed:", err.Error())
	}

}

func (peer *Peer) String() string {
	return peer.IpAddress + ":" + strconv.FormatUint(uint64(peer.Port), 10)
}
