package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

/// ///

type Router struct {
	multiplexer    *http.ServeMux
	mappedHandlers map[string]map[HttpMethod]ActionFunc
	controllers    []Controller
}

/// ///

func (r *Router) IfIndexNotDefinedAdd(control Controller) *Router {
	if _, exist := r.mappedHandlers["/"]; !exist {
		r.AddControl(
			control,
		)
	}

	return r
}

/// ///

func (r *Router) AddControl(controller Controller) *Router {
	if r.mappedHandlers == nil {
		r.mappedHandlers = make(map[string]map[HttpMethod]ActionFunc)
	}

	for _, endpoint := range controller.Endpoints().endpoints {
		if endpoint.Path == "/" {
			endpoint.Path = ""
		}

		fullpath := controller.Basepath() + endpoint.Path

		if _, exists := r.mappedHandlers[fullpath]; !exists {
			r.mappedHandlers[fullpath] = make(map[HttpMethod]ActionFunc)
		}

		for method := range r.mappedHandlers[fullpath] {
			if method == endpoint.Method {
				panic("Method " + string(endpoint.Method) + " already registered at " + fullpath)
			}
		}

		r.mappedHandlers[fullpath][endpoint.Method] = endpoint.Action
	}

	r.controllers = append(r.controllers, controller)

	return r
}

/// ///

func (r *Router) Serve() *http.ServeMux {
	r.multiplexer = http.NewServeMux()

	for path, handlers := range r.mappedHandlers {
		var pathHandlers *map[HttpMethod]*ActionFunc = new(map[HttpMethod]*ActionFunc)
		*pathHandlers = make(map[HttpMethod]*ActionFunc)

		fmt.Printf("%v [", path)

		for method := range handlers {
			(*pathHandlers)[method] = new(ActionFunc)
			*(*pathHandlers)[method] = handlers[method]

			fmt.Printf(" %v", method)
		}

		fmt.Printf(" ]\n")

		execute := func(writer http.ResponseWriter, request *http.Request) {
			var response Response

			if action, mappedHandler := (*pathHandlers)[HttpMethod(request.Method)]; mappedHandler {
				response = (*action)(NewRequest(&writer, request))
			} else {
				response = ErrorResponse().
					Message("Method not defined for path, can't handle request")
			}

			if response != nil {
				switch response.getType() {
				case responseHtml:
					fallthrough
				case responseCss:
					fallthrough
				case responseJavaScript:
					fallthrough
				case responseApplication:
					fallthrough
				case responseFile:
					file, _ := os.Open(response.getMeta(responseMetaFileRelativePath).(string))

					http.ServeContent(
						writer, request, response.getMeta(responseMetaFileRelativePath).(string), time.Time{}, file,
					)
				default:
					for key, data := range response.getHeader() {
						for _, val := range data {
							writer.Header().Add(key, val)
						}
					}
					writer.WriteHeader(int(response.getStatus()))
					writer.Write(response.getResponse())
				}
			}
		}

		r.multiplexer.HandleFunc(path, execute)
		r.multiplexer.HandleFunc(path+"/", execute)
	}

	return r.multiplexer
}

/// ///
