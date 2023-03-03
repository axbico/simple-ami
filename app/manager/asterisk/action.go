package asterisk

import "fmt"

/// ///

type ActionType string

const (
	ACTION_BRIDGE_LIST          ActionType = "BridgeList"
	ACTION_SHOW_CHANNELS        ActionType = "CoreShowChannels"
	ACTION_EXTENSION_STATE_LIST ActionType = "ExtensionStateList"
	ACTION_LOGIN                ActionType = "Login"
	ACTION_SHOW_CONTACTS        ActionType = "PJSIPShowContacts"
	ACTION_SHOW_ENDPOINTS       ActionType = "PJSIPShowEndpoints"
)

/// ///

type action struct {
	packetEntity
}

func NewAction() *action {
	var act *action = new(action)
	act.raw = ""
	act.status = ""
	act.headers = make(map[string]string)
	return act
}

func (ac *action) Action(value ActionType) *action {
	ac.SetStatus(string(value), ACTION)
	return ac
}

func (ac *action) Header(key string, value string) *action {
	ac.SetHeader(key, value)
	return ac
}

func (ac *action) Raw() string {
	return ac.raw
}

func (ac *action) Packet() string {
	if ac.status == "" {
		fmt.Println("No action defined")
		return ""
	}

	return ac.raw + "\r\n"
}

/// ///
