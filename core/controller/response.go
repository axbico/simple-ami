package controller

import (
	"fmt"
	"net/http"
)

/// ///

type (
	responseType uint16
	statusCode   int
)

type Response interface {
	getResponse() []byte
	getStatus() statusCode
	getHeader() http.Header
	getMeta(string) interface{}
	getType() responseType
}

/// ///

const (
	responseDefault responseType = 1 << iota
	responseError
	responseRaw
	responseJson
	responseFile
	responseHtml
	responseCss
	responseJavaScript
	responseApplication
)

const (
	responseMetaRequest          string = "meta-request"
	responseMetaFileRelativePath string = "meta-file-relative-path"
)

const (
	contentTypeText        = "text/plain"
	contentTypeHtml        = "text/html"
	contentTypeCss         = "text/css"
	contentTypeJavascript  = "application/javascript"
	contentTypeJson        = "application/json"
	contentTypeOctetStream = "application/octet-stream"
)

/// ///

type response struct {
	status       statusCode
	responseType responseType
	meta         map[string]interface{}
	header       http.Header
	rawContent   []byte
	message      string
}

func (r *response) Status(status statusCode) *response {
	r.status = status
	return r
}

func (r *response) Message(msg string) *response {
	r.message = r.message + msg
	r.newRawContentFromMessage()
	return r
}

func (r *response) Filepath(path string) *response {
	r.meta[responseMetaFileRelativePath] = path
	return r
}

func (r *response) Request(request *http.Request) *response {
	r.meta[responseMetaRequest] = request
	return r
}

/// ///

func newResponse() *response {
	return &response{
		responseType: responseDefault,
		message:      "",
		status:       http.StatusOK,
		rawContent:   []byte{},
		header:       http.Header{},
		meta:         make(map[string]interface{}),
	}
}

func (r *response) newRawContentFromMessage() *response {
	r.rawContent = []byte(r.message)
	return r
}

func (r *response) getResponse() []byte {
	r.Message(fmt.Sprint(r.responseType))
	return r.rawContent
}

func (r *response) getStatus() statusCode {
	return r.status
}

func (r *response) getHeader() http.Header {
	return r.header
}

func (r *response) getMeta(key string) interface{} {
	if value, set := r.meta[key]; set {
		return value
	}

	panic("Response meta value for (" + key + ") doesn't exist")
}

func (r *response) getType() responseType {
	return r.responseType
}

func (r *response) setHeaderContentType(ct string) *response {
	r.header.Add("content-type", ct)
	return r
}

/// ///

func RawResponse() *response {
	var r *response = newResponse()
	r.responseType = responseRaw
	r.setHeaderContentType(contentTypeText)
	return r
}

func ErrorResponse() *response {
	var r *response = newResponse()
	r.responseType = responseError
	r.status = http.StatusBadRequest
	r.setHeaderContentType(contentTypeText)
	return r
}

func JsonResponse() *response {
	var r *response = newResponse()
	r.responseType = responseJson
	r.setHeaderContentType(contentTypeJson)
	return r
}

func FileResponse() *response {
	var r *response = newResponse()
	r.responseType = responseFile
	r.setHeaderContentType(contentTypeOctetStream)
	return r
}

func HtmlResponse() *response {
	var r *response = newResponse()
	r.responseType = responseHtml
	r.setHeaderContentType(contentTypeHtml)
	return r
}

func CssResponse() *response {
	var r *response = newResponse()
	r.responseType = responseCss
	r.setHeaderContentType(contentTypeCss)
	return r
}

func JavascriptResponse() *response {
	var r *response = newResponse()
	r.responseType = responseJavaScript
	r.setHeaderContentType(contentTypeJavascript)
	return r
}

/// ///

func ApplicationResponse(request *http.Request) *response {
	var r *response = newResponse()
	r.responseType = responseApplication
	r.setHeaderContentType(contentTypeOctetStream)
	return r
}

/// ///
