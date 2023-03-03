package asterisk

import "fmt"

/// ///

type EventType string

const (
	EVENT_BRIDGE_CREATE                 EventType = "BridgeCreate"
	EVENT_BRIDGE_DESTROY                EventType = "BridgeDestroy"
	EVENT_BRIDGE_ENTER                  EventType = "BridgeEnter"
	EVENT_BRIDGE_LEAVE                  EventType = "BridgeLeave"
	EVENT_BRIDGE_LIST_COMPLETE          EventType = "BridgeListComplete"
	EVENT_CONTACT_LIST_COMPLETE         EventType = "ContactListComplete"
	EVENT_CONTACT_STATUS                EventType = "ContactStatus"
	EVENT_CORE_SHOW_CHANNEL             EventType = "CoreShowChannel"
	EVENT_ENDPOINT_DETAIL               EventType = "EndpointDetail"
	EVENT_ENDPOINT_DETAIL_COMPLETE      EventType = "EndpointDetailComplete"
	EVENT_ENDPOINT_LIST                 EventType = "EndpointList"
	EVENT_ENDPOINT_LIST_COMPLETE        EventType = "EndpointListComplete"
	EVENT_EXTENSION_STATUS              EventType = "ExtensionStatus"
	EVENT_EXTENSION_STATE_LIST_COMPLETE EventType = "ExtensionStateListComplete"
	EVENT_HANGUP                        EventType = "Hangup"
	EVENT_NEWSTATE                      EventType = "Newstate"
	EVENT_PEER_STATUS                   EventType = "PeerStatus"
	EVENT_RELOAD                        EventType = "Reload"
)

/// ///

type event struct {
	packetEntity
}

func Event(pe *packetEntity) *event {
	var e *event = new(event)
	e.packetEntity = *pe
	return e
}

func (e *event) Status() string {
	return e.GetStatus()
}

func (e *event) Header(key string) string {
	return e.GetHeader(key)
}

func (e *event) Raw() string {
	return e.raw
}

func PrintEvent(eventChannel chan *event, watch string) {
	for {
		e := <-eventChannel

		if e.status == watch {
			fmt.Println("== EVENT ==")
			fmt.Println(e.raw)
		}
	}
}

/// ///
