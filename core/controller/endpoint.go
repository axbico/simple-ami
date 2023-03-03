package controller

import (
	"fmt"
	"net/http"
	"strings"
)

/// ///

type HttpMethod string

type ActionFunc func(*Request) Response

type Controller interface {
	Basepath() string
	Endpoints() ControllerEndpoints
}

/// ///

type Endpoint struct {
	Action ActionFunc
	Method HttpMethod
	Path   string
}

type ControllerEndpoints struct {
	endpoints []Endpoint
}

/// ///

func (ce *ControllerEndpoints) AddEndpoint(path string, method HttpMethod, action ActionFunc) *ControllerEndpoints {
	method = HttpMethod(strings.ToUpper(string(method)))

	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodTrace, http.MethodConnect:
	default:
		fmt.Println("registering method [" + string(method) + "] is not allowed")
	}

	ce.endpoints = append(
		ce.endpoints,
		Endpoint{
			Path:   path,
			Method: method,
			Action: action,
		},
	)

	return ce
}

func (ce *ControllerEndpoints) GroupExpand(c ...Controller) *ControllerEndpoints {
	if len(c) > 0 {
		for _, controller := range c {
			for _, endpoint := range controller.Endpoints().endpoints {
				ce.AddEndpoint(
					controller.Basepath()+endpoint.Path,
					endpoint.Method,
					endpoint.Action,
				)
			}
		}
	}

	return ce
}

/// ///

func (ce *ControllerEndpoints) GET(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodGet, a)
}

func (ce *ControllerEndpoints) POST(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodPost, a)
}

func (ce *ControllerEndpoints) PUT(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodPut, a)
}

func (ce *ControllerEndpoints) DELETE(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodDelete, a)
}

func (ce *ControllerEndpoints) PATCH(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodPatch, a)
}

func (ce *ControllerEndpoints) HEAD(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodHead, a)
}

func (ce *ControllerEndpoints) OPTIONS(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodOptions, a)
}

func (ce *ControllerEndpoints) TRACE(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodTrace, a)
}

func (ce *ControllerEndpoints) CONNECT(p string, a ActionFunc) *ControllerEndpoints {
	return ce.AddEndpoint(p, http.MethodConnect, a)
}

/// ///
