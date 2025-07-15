package torrent

import (
	"encoding/binary"
)

type MessageType int
const (
	KEEP_ALIVE MessageType = -1
	CHOKE MessageType = 0
	UNCHOKE MessageType = 1
	INTERESTED MessageType = 2	
	NOT_INTERESTED MessageType = 3
	HAVE MessageType = 4
	BITFIELD MessageType = 5
	REQUEST MessageType = 6
	PIECE MessageType = 7
	CANCEL MessageType = 8
	PORT MessageType = 9
)

type Message struct {
	messageType MessageType
	payload []byte
}

func parseMessage(rawMessage []byte) Message {
	if len(rawMessage) < 4 {
	}
	length := int(binary.BigEndian.Uint16(rawMessage[:4]))	
	if length == 0 {
		// keep alive msg, reset read deadline to two mins	
		return Message{
			messageType: KEEP_ALIVE,
		}
	}
	messageId := int(rawMessage[4])
	var payload []byte
	if length > 0 {
		payload = rawMessage[5:]	
	}

	return Message{
		messageType: MessageType(messageId),
		payload: payload,
	}
}
