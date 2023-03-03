package feedgroup

import (
	"pbx/amitask/app/manager/asterisk"
	"pbx/amitask/app/monitor"
	"pbx/amitask/core/controller"
	"pbx/amitask/core/misc"
)

/// ///

type DashboardControl struct {
	//
}

func (control *DashboardControl) Basepath() string {
	return "/dashboard"
}

func (control *DashboardControl) Endpoints() controller.ControllerEndpoints {
	return *(new(controller.ControllerEndpoints)).
		GET(
			"/connect", Dashboard,
		)
}

/// ///

func Dashboard(r *controller.Request) controller.Response {

	var websocket *misc.Websocket = misc.WebsocketUpgradeConnection(r)

	if websocket == nil {
		return controller.ErrorResponse()
	}

	var ami *asterisk.AMI = asterisk.ConnectAMI("tcp", ":5038")
	defer ami.Disconnect()

	go ami.Monitor(monitor.DashboardMonitor())

	go func(ami *asterisk.AMI, websocket *misc.Websocket) {
		for {
			if packets := <-ami.UnifiedOutput(); packets != nil {
				for _, output := range packets {
					if output != nil {
						websocket.Send([]byte(output.(string)))
					}
				}
			}
		}
	}(ami, websocket)

	<-websocket.WaitTerminatedConnectionChannel()

	return nil
}

/// ///
