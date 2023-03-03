package asterisk

import "fmt"

/// ///

type response struct {
	packetEntity
}

func Response(pe *packetEntity) *response {
	var r *response = new(response)
	r.packetEntity = *pe
	return r
}

func (r *response) Response() string {
	return r.GetStatus()
}

func (r *response) Header(key string) string {
	return r.GetHeader(key)
}

func (r *response) Raw() string {
	return r.raw
}

func PrintResponse(responseChannel chan *response) {
	for {
		r := <-responseChannel

		fmt.Println("== RESPONSE ==")
		fmt.Println(r.raw)
	}
}

/// ///
