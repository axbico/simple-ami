package monitor

import (
	"fmt"
	"pbx/amitask/app/manager/asterisk"
)

/// ///

const (
	DASHBOARD_ALL_USERS         asterisk.MonitorTag = "ALL_USERS"
	DASHBOARD_ALL_USERSverbose  asterisk.MonitorTag = "ALL_USERSverbose"
	DASHBOARD_ACTIVE_USERS      asterisk.MonitorTag = "ACTIVE_USERS"
	DASHBOARD_ACTIVE_CALLS      asterisk.MonitorTag = "ACTIVE_CALLS"
	DASHBOARD_RECENT_ACTIVITIES asterisk.MonitorTag = "RECENT_ACTIVITIES"
	DASHBOARD_CALL_ACTIVITIES   asterisk.MonitorTag = "CALL_ACTIVITIES"
)

/// ///

func DashboardMonitor() *asterisk.Monitor {
	var dashboardMonitor *asterisk.Monitor = asterisk.NewMonitor()

	// get total number of endpoints
	dashboardMonitor.
		AddMonitor(
			DASHBOARD_ALL_USERS,
			asterisk.EVENT_ENDPOINT_LIST_COMPLETE,
			allUsersCompleteCount,
		).
		AddTriggerAction(asterisk.NewAction().Action(asterisk.ACTION_SHOW_ENDPOINTS))

	// get number of endpoints with active contact
	dashboardMonitor.
		AddMonitor(
			DASHBOARD_ACTIVE_USERS,
			asterisk.EVENT_CONTACT_LIST_COMPLETE,
			activeUsersCompleteCount,
		).
		AddTriggerAction(asterisk.NewAction().Action(asterisk.ACTION_SHOW_CONTACTS))

	// monitor for changes in number of endpoints with active contact
	dashboardMonitor.
		AddMonitor(
			DASHBOARD_ACTIVE_USERS,
			asterisk.EVENT_CONTACT_STATUS,
			activeUsersStatus,
		)

	// monitor for logging basic status of users
	dashboardMonitor.
		AddMonitor( // online/offline contact
			DASHBOARD_RECENT_ACTIVITIES,
			asterisk.EVENT_CONTACT_STATUS,
			recentActivitiesContactStatus,
		).
		AddMonitor( // active calls between extensions
			DASHBOARD_RECENT_ACTIVITIES,
			asterisk.EVENT_NEWSTATE,
			recentActivitiesContactStatus,
		)

	// monitor for number of active calls
	dashboardMonitor.
		AddMonitor( // monitor call start
			DASHBOARD_ACTIVE_CALLS,
			asterisk.EVENT_NEWSTATE,
			activeOutgoingCallsUpdate,
		).
		AddMonitor( // monitor call end
			DASHBOARD_ACTIVE_CALLS,
			asterisk.EVENT_HANGUP,
			activeOutgoingCallsUpdate,
		).
		AddMonitor( // initial number of active calls
			DASHBOARD_ACTIVE_CALLS,
			asterisk.EVENT_CORE_SHOW_CHANNEL,
			activeOutgoingCallsUpdate,
		).
		AddTriggerAction(asterisk.NewAction().Action(asterisk.ACTION_SHOW_CHANNELS))

	// monitor details for logging channel status of active calls
	dashboardMonitor.
		AddMonitor(
			DASHBOARD_CALL_ACTIVITIES,
			asterisk.EVENT_NEWSTATE,
			callOutgoingActivitiesStatus,
		).
		AddMonitor(
			DASHBOARD_CALL_ACTIVITIES,
			asterisk.EVENT_HANGUP,
			callOutgoingActivitiesStatus,
		)

	dashboardMonitor.AddLogin("username", "secret")

	return dashboardMonitor
}

/// ///

func allUsersCompleteCount(packet asterisk.Packet) interface{} {
	return string(DASHBOARD_ALL_USERS) + "~" + packet.GetHeader("ListItems")
}

/// ///

func activeUsersCompleteCount(packet asterisk.Packet) interface{} {
	return string(DASHBOARD_ACTIVE_USERS) + "~" + packet.GetHeader("ListItems")
}

/// ///

func activeUsersStatus(packet asterisk.Packet) interface{} {
	var count int = 1
	status := packet.GetHeader("ContactStatus")

	switch status {
	case "Removed":
		count *= -1
	case "NonQualified":
	default:
		return nil
	}

	return "ACTIVE_USERS_ALTER~" + fmt.Sprint(count)
}

/// ///

func recentActivitiesContactStatus(packet asterisk.Packet) interface{} {
	if packet.GetStatus() == string(asterisk.EVENT_NEWSTATE) {
		if packet.GetHeader("Uniqueid") == packet.GetHeader("Linkedid") && packet.GetHeader("ChannelStateDesc") == "Up" {
			return "RECENT_ACTIVITIES~" + packet.GetHeader("CallerIDNum") + " <i class=\"phone icon\" style=\"margin: 0 10px;\"></i> " + packet.GetHeader("Exten")
		} else {
			return nil
		}
	}

	var status string = packet.GetHeader("ContactStatus")
	switch status {
	case "Removed":
		status = "disconnected"
	case "NonQualified":
		status = "is online"
	default:
		status = "status unknown"
	}
	return "RECENT_ACTIVITIES~User " + packet.GetHeader("EndpointName") + " " + status
}

/// ///

func activeOutgoingCallsUpdate(packet asterisk.Packet) interface{} {
	if packet.GetHeader("ChannelStateDesc") != "Up" || packet.GetHeader("Uniqueid") != packet.GetHeader("Linkedid") {
		return nil
	}

	var count int = 1
	switch packet.GetStatus() {
	case string(asterisk.EVENT_HANGUP):
		count *= -1
	default:
	}

	return "ACTIVE_CALLS_ALTER~" + fmt.Sprint(count)
}

/// ///

func callOutgoingActivitiesStatus(packet asterisk.Packet) interface{} {
	var message string = ""

	var callerId string = packet.GetHeader("CallerIDNum")
	var exten string = packet.GetHeader("Exten")
	var callId string = packet.GetHeader("Uniqueid")

	if packet.GetHeader("Uniqueid") != packet.GetHeader("Linkedid") {
		callerId = packet.GetHeader("ConnectedLineNum")
		exten = packet.GetHeader("CallerIDNum")
		callId = packet.GetHeader("Linkedid")
	}

	switch packet.GetStatus() {
	case string(asterisk.EVENT_NEWSTATE):
		switch packet.GetHeader("ChannelStateDesc") {
		case "Ringing":
			message += callerId + " called extension " + exten + "~" + callId
		case "Up":
			if packet.GetHeader("Uniqueid") != packet.GetHeader("Linkedid") {
				message += exten + " accepted call from " + callerId + "~" + callId
			} else {
				message += callerId + " on call with " + exten + "~" + callId
			}
		case "Busy":
			message += callerId + " called busy extension " + exten
		default:
			return nil
		}
	case string(asterisk.EVENT_HANGUP):
		if packet.GetHeader("Uniqueid") != packet.GetHeader("Linkedid") {
			hang := callerId
			if packet.GetHeader("Cause") == "17" || packet.GetHeader("Cause") == "16" {
				hang = exten
			}
			message += "Hangup action on " + callId + " by extension " + hang
		} else {
			if packet.GetHeader("ChannelStateDesc") == "Busy" || packet.GetHeader("ChannelStateDesc") == "Ringing" {
				return nil
			}
			message += "Call " + callId + " from user " + callerId + " to " + exten + " completed"
		}
	default:
	}

	return "CALL_ACTIVITIES~" + message
}

/// ///
