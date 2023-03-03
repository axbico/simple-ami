package monitor

import (
	"pbx/amitask/app/manager/asterisk"
)

/// ///

const (
	USERS_ALL_USERS asterisk.MonitorTag = "ALL_USERS"
)

/// ///

func UsersMonitor() *asterisk.Monitor {
	var usersMonitor *asterisk.Monitor = asterisk.NewMonitor()

	usersMonitor.
		AddMonitor( // get all users
			USERS_ALL_USERS,
			asterisk.EVENT_ENDPOINT_LIST,
			listEndpoint,
		).
		AddMonitor( // monitor changes for endpoints active contact
			USERS_ALL_USERS,
			asterisk.EVENT_CONTACT_STATUS,
			contactStatusUpdate,
		).
		AddTriggerAction(asterisk.NewAction().Action(asterisk.ACTION_SHOW_ENDPOINTS))

	usersMonitor.AddLogin("username", "secret")

	return usersMonitor
}

/// ///

func listEndpoint(packet asterisk.Packet) interface{} {
	user := "USERS_ALL_USERS~" + packet.GetHeader("ObjectName")
	if packet.GetHeader("Contacts") != "" {
		user += "~" + packet.GetHeader("Contacts")
	}
	return user
}

/// ///

func contactStatusUpdate(packet asterisk.Packet) interface{} {
	var status string = packet.GetHeader("ContactStatus")
	switch status {
	case "Removed":
		status = "REMOVE" + "~" + packet.GetHeader("EndpointName")
	case "NonQualified":
		status = "ADD" + "~" + packet.GetHeader("EndpointName") + "~" + packet.GetHeader("URI")
	default:
		status = "UNKNOWN" + "~" + packet.GetHeader("EndpointName") + "~" + packet.GetHeader("URI")
	}
	return "USERS_CONTACT_" + status
}

/// ///
