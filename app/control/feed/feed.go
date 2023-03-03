package control

import (
	feedgroup "pbx/amitask/app/control/feed/group"
	"pbx/amitask/core/controller"
)

/// ///

type FeedControl struct {
	// just groups controllers
}

func (control *FeedControl) Basepath() string {
	return "/api"
}

func (control *FeedControl) Endpoints() controller.ControllerEndpoints {
	return *(new(controller.ControllerEndpoints)).
		GroupExpand(
			new(feedgroup.UserControl),
			new(feedgroup.DashboardControl),
		)
}

/// ///
