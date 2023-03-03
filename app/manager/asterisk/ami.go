package asterisk

import (
	"bufio"
	"fmt"
	"net"
)

/// ///

type AMI struct {
	connection   *net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	lastAction   *action
	lastResponse *response
	lastEvent    *event

	monitor *Monitor

	responseChannel chan *response
	eventChannel    chan *event
	unifiedOutput   chan []interface{}
}

/// ///

func ConnectAMI(network, address string) *AMI {
	var ami *AMI = new(AMI)

	conn, err := net.Dial(network, address)

	if err != nil {
		panic("Problem connecting to Asterisk")
	}

	ami.connection = &conn
	ami.reader = bufio.NewReader(*ami.connection)
	ami.writer = bufio.NewWriter(*ami.connection)
	ami.lastAction = nil
	ami.lastResponse = nil
	ami.lastEvent = nil

	ami.responseChannel = make(chan *response)
	ami.eventChannel = make(chan *event)
	ami.unifiedOutput = make(chan []interface{})

	return ami
}

/// ///

func (ami *AMI) Response() string {
	return ami.lastResponse.GetRaw()
}

func (ami *AMI) Event() string {
	return ami.lastEvent.GetRaw()
}

func (ami *AMI) ResponseChannel() chan *response {
	return ami.responseChannel
}

func (ami *AMI) EventChannel() chan *event {
	return ami.eventChannel
}

func (ami *AMI) UnifiedOutput() chan []interface{} {
	return ami.unifiedOutput
}

/// ///

func (ami *AMI) Disconnect() *AMI {
	(*ami.connection).Close()

	return ami
}

func (ami *AMI) Send(action *action) *AMI {
	ami.lastAction = action

	if _, err := ami.writer.WriteString(ami.lastAction.Packet()); err != nil {
		fmt.Println("Failed to write in buffer")
		ami.lastAction = nil
		ami.lastResponse = nil

		return ami
	}

	ami.writer.Flush()

	return ami
}

/// ///

func (ami *AMI) Monitor(m *Monitor) {
	ami.monitor = m

	go ami.Subscribe(
		ami.monitor.handleResponse,
		ami.monitor.handleEvent,
	)

	ami.Send(ami.monitor.login)

	for _, action := range ami.monitor.triggerActions {
		ami.Send(action)
	}

	<-make(chan interface{})
}

/// ///

func (ami *AMI) Subscribe(responseHandler func(Packet) []interface{}, eventHandler func(Packet) []interface{}) {

	for {
		switch packet, pType := readPacketFromBufferedReader(ami.reader); pType {
		case RESPONSE:
			ami.unifiedOutput <- responseHandler(packet)
		case EVENT:
			ami.unifiedOutput <- eventHandler(packet)
		}
	}
}

/// ///
