package server

import (
	controlApplication "pbx/amitask/app/control/application"
	controlFeed "pbx/amitask/app/control/feed"
	"pbx/amitask/core/controller"
	"pbx/amitask/core/misc"
)

/// ///

func AllRouter() *controller.Router {
	return new(controller.Router).
		AddControl(
			new(controlApplication.ApplicationControl),
		).
		AddControl(
			new(controlFeed.FeedControl),
		).
		IfIndexNotDefinedAdd(
			new(misc.IndexControl),
		)
}

/// ///
