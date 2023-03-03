package asterisk

import (
	"bufio"
	"strings"
)

/// ///

type packetType string

const (
	ACTION      packetType = "Action"
	RESPONSE    packetType = "Response"
	EVENT       packetType = "Event"
	UNKNOW_TYPE packetType = "UnknownType"
)

/// ///

type packet interface {
	SetHeader(string, string) packet
	GetHeader(string) string
	SetStatus(string, packetType) packet
	GetStatus() string
	SetRaw(string) packet
	GetRaw() string
	Raw() string
}

type Packet interface {
	SetHeader(string, string) packet
	GetHeader(string) string
	SetStatus(string, packetType) packet
	GetStatus() string
	SetRaw(string) packet
	GetRaw() string
	Raw() string
}

/// ///

type packetEntity struct {
	raw     string
	headers map[string]string
	status  string
}

func NewPacketEntity() *packetEntity {
	var pe *packetEntity = new(packetEntity)
	pe.raw = ""
	pe.status = ""
	pe.headers = make(map[string]string)
	return pe
}

func (pe *packetEntity) SetHeader(header string, value string) packet {
	pe.headers[header] = value
	pe.raw += newHeaderLine(header, value)

	return pe
}

func (pe *packetEntity) GetHeader(header string) string {
	if value, ok := pe.headers[header]; ok {
		return value
	}

	return ""
}

func (pe *packetEntity) SetStatus(value string, t packetType) packet {
	pe.status = value
	pe.raw = newHeaderLine(string(t), pe.status) + pe.raw

	return pe
}

func (pe *packetEntity) GetStatus() string {
	return pe.status
}

func (pe *packetEntity) SetRaw(r string) packet {
	pe.raw = r
	return pe
}

func (pe *packetEntity) GetRaw() string {
	return pe.raw
}

func (pe *packetEntity) Raw() string {
	return pe.raw
}

/// ///

func newHeaderLine(key string, value string) string {
	return key + ": " + value + "\r\n"
}

func readPacketFromBufferedReader(reader *bufio.Reader) (Packet, packetType) {
	var pe *packetEntity = NewPacketEntity()

	for {
		line, err := reader.ReadString('\n')
		if line == "\r\n" || err != nil {
			break
		}
		pe.raw += line
		if ix := strings.Index(line, ":"); ix != -1 {
			pe.headers[strings.TrimSpace(line[:ix])] = strings.TrimSpace(line[ix+1:])
		}
	}

	// first check for event cause RESPONSE can be present as event header
	if status, ok := pe.headers[string(EVENT)]; ok {
		pe.status = status
		return Event(pe), EVENT
	} else if status, ok := pe.headers[string(RESPONSE)]; ok {
		pe.status = status
		return Response(pe), RESPONSE
	}

	return pe, UNKNOW_TYPE
}

/// ///
