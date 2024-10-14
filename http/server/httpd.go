package server

import (
	"net/http"
)

func StartTestHttpServer(server *http.Server) *http.Server {
	if server.TLSConfig == nil {
		go server.ListenAndServe()
	} else {
		go server.ListenAndServeTLS("", "")
	}
	return server
}
