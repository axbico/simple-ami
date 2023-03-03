package controller

import "net/http"

/// ///

type Request struct {
	Writer  *http.ResponseWriter
	Request *http.Request
}

func NewRequest(writer *http.ResponseWriter, request *http.Request) *Request {
	var r *Request = new(Request)
	r.Writer = writer
	r.Request = request

	return r
}

/// ///
