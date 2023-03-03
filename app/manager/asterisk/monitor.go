package asterisk

/// ///

type MonitorTag string

/// ///

type Monitor struct {
	monitors       map[EventType]map[MonitorTag]func(Packet) interface{}
	triggerActions []*action
	login          *action
}

func NewMonitor() *Monitor {
	var monitor *Monitor = new(Monitor)

	monitor.monitors = make(map[EventType]map[MonitorTag]func(Packet) interface{})
	monitor.triggerActions = []*action{}
	monitor.login = nil

	return monitor
}

/// ///

/*
by having monitor tag in structure, monitor allows handling multiple times of one AMI event from one connection
*/
func (monitor *Monitor) AddMonitor(tag MonitorTag, event EventType, handler func(Packet) interface{}) *Monitor {
	if _, ok := monitor.monitors[event]; !ok {
		monitor.monitors[event] = make(map[MonitorTag]func(Packet) interface{})
	}

	monitor.monitors[event][tag] = handler

	return monitor
}

func (monitor *Monitor) AddLogin(username string, secret string) *Monitor {
	monitor.login = NewAction().Action(ACTION_LOGIN).
		Header("Username", username).
		Header("Secret", secret)

	return monitor
}

func (monitor *Monitor) AddTriggerAction(a *action) *Monitor {
	monitor.triggerActions = append(monitor.triggerActions, a)
	return monitor
}

/// ///

func (monitor *Monitor) handleResponse(p Packet) []interface{} {
	return nil
}

/// ///

func (monitor *Monitor) handleEvent(p Packet) []interface{} {
	// handle triggered event for each MonitorTag
	if tags, ok := monitor.monitors[EventType(p.GetStatus())]; ok {
		var output []interface{} = []interface{}{}
		for _, handler := range tags {
			output = append(output, handler(p))
		}
		return output
	}

	return nil
}

/// ///
