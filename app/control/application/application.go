package control

import (
	"net/http"
	"pbx/amitask/core/controller"
	"pbx/amitask/core/misc"
)

/// ///

type ApplicationControl struct {
	// controller for serving application/frontend files at root directory of process from /application
}

func (control *ApplicationControl) Basepath() string {
	return "/app"
}

func (control *ApplicationControl) Endpoints() controller.ControllerEndpoints {
	return *(new(controller.ControllerEndpoints)).
		GET(
			"/", Index,
		)
}

/// ///

func Index(request *controller.Request) controller.Response {
	if request.Request.URL.Path == new(ApplicationControl).Basepath()+"/" {
		request.Request.URL.Path += "app.html"
	}
	return PropagatePublicAccess(request.Request, true)
}

/// ///

func PropagatePublicAccess(request *http.Request, allowPublicAccess ...bool) controller.Response {
	if len(allowPublicAccess) > 0 && allowPublicAccess[0] {
		request.URL.Path = request.URL.Path[len(new(ApplicationControl).Basepath())+1:]

		if filepath, allowed := misc.AllowedFileApplication(request); allowed {
			return controller.ApplicationResponse(request).
				Filepath(filepath)
		}
	}

	return controller.ApplicationResponse(request).Message("Application resource " + request.URL.RequestURI() + " not found")
}

/// ///
