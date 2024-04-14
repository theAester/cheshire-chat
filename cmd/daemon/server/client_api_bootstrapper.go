package server

import (
  "fmt"
  "net/http"

	//"github.com/theAester/cheshire-chat/cmd/daemon/util"
)

func SetupClientApi (server *Server, mux *http.ServeMux) {
  // Here we add end points for client API
  mux.HandleFunc("/v1/test", func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from daemon")
  })
}
