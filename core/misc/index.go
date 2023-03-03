package misc

import (
	"net/http"
	"pbx/amitask/core/controller"
)

/// ///

type IndexControl struct {
	//
}

func (control *IndexControl) Basepath() string {
	return "/"
}

func (control *IndexControl) Endpoints() controller.ControllerEndpoints {
	return *(new(controller.ControllerEndpoints)).
		GET(
			"/",
			func(request *controller.Request) controller.Response {
				if request.Request.URL.Path == "/" {
					request.Request.URL.Path += "index.html"
				}
				return NotFoundHandler(request.Request, true)
			},
		).
		POST(
			"/",
			func(request *controller.Request) controller.Response {
				return controller.RawResponse().Message("JUST FOR TEST POST WORKS")
			},
		)
}

/// ///

func NotFoundHandler(request *http.Request, allowPublicAccess ...bool) controller.Response {
	if len(allowPublicAccess) > 0 && allowPublicAccess[0] {
		if filepath, allowed := AllowedFilePublic(request); allowed {
			return controller.FileResponse().
				Filepath(filepath).
				Request(request)
		}
	}

	return controller.ErrorResponse().Message("Not found request URI " + request.URL.RequestURI())
}

/// ///
