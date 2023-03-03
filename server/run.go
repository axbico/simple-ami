package server

import (
	"fmt"
	"log"
	"net/http"
	"pbx/amitask/core/controller"
)

/// ///

func Run(domain string, port int32, router *controller.Router) {
	log.Fatal(
		http.ListenAndServe(
			domain+":"+fmt.Sprintf("%v", port),
			router.Serve(),
		),
	)
}

/// ///
