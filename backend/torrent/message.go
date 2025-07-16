package torrent

import (
	"encoding/binary"
	"fmt"
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

func (message *Message) messageTypeToString() string {
	switch message.messageType {
	case KEEP_ALIVE:
		return fmt.Sprint("KEEP_ALIVE")
	case CHOKE:
		return fmt.Sprint("CHOKE")
	case UNCHOKE:
		return fmt.Sprint("UNCHOKE")
	case INTERESTED:
		return fmt.Sprint("INTERESTED")
	case NOT_INTERESTED:
		return fmt.Sprint("NOT_INTERESTED")
	case HAVE:
		return fmt.Sprint("HAVE")
	case BITFIELD:
		return fmt.Sprint("BITFIELD")
	case REQUEST:
		return fmt.Sprint("REQUEST")
	case PIECE:
		return fmt.Sprint("PIECE")
	case CANCEL:
		return fmt.Sprint("CANCEL")
	case PORT:
		return fmt.Sprint("PORT")
	}
	return ""
}

func (message *Message) String() string {
	return fmt.Sprintf("Message Type: %s, Message Payload: %s\n", message.messageTypeToString(), string(message.payload))
}
