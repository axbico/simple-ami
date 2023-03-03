package misc

import (
	"pbx/amitask/core/controller"

	"github.com/gorilla/websocket"
)

/// ///

type Websocket struct {
	connection                      *websocket.Conn
	waitTerminatedConnectionChannel chan bool
}

func WebsocketUpgradeConnection(r *controller.Request) *Websocket {
	var ws *Websocket = new(Websocket)
	ws.waitTerminatedConnectionChannel = make(chan bool)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
	}

	connection, err := upgrader.Upgrade(*r.Writer, r.Request, nil)
	if err != nil {
		return nil
	}

	ws.connection = connection

	ws.listenForTermination()

	return ws
}

/// ///

func (ws *Websocket) CloseConnection() {
	ws.connection.Close()
	ws.waitTerminatedConnectionChannel <- true
}

/// ///

func (ws *Websocket) Send(output []byte) {
	if err := ws.connection.WriteMessage(websocket.TextMessage, output); err != nil {
		ws.CloseConnection()
	}
}

/// ///

func (ws *Websocket) listenForTermination() {
	go func() {
		for {
			if x, _, _ := ws.connection.ReadMessage(); x == -1 {
				ws.CloseConnection()
			}
		}
	}()
}

/// ///

func (ws *Websocket) WaitTerminatedConnectionChannel() chan bool {
	return ws.waitTerminatedConnectionChannel
}

/// ///
