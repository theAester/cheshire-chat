package server

import (
	//"github.com/theAester/cheshire-chat"
	"fmt"
  "os"
  "net/http"

	"github.com/theAester/cheshire-chat/cmd/daemon/util"
	//"github.com/theAester/cheshire-chat/cmd/daemon/util/negotiator"
	//"github.com/theAester/cheshire-chat/cmd/daemon/util/client_api"
	"github.com/theAester/cheshire-chat/share/security"
)

type Server struct {
  ClientServer *http.Server
  ClientSocketPath string
  //Negotiator negotiator.Negotiator
  //ConMan util.ConnectionManager
}



func New(config *util.Config) (*Server, error) {
  err := security.InitPGP(config)
  if err != nil {
    return nil, fmt.Errorf("Cannot initalize PGP: %v", err)
  }

  //err := negotiator.Init(config)
  //if err != nil {
    //return nil, fmt.Errorf("Cannot create negotiator: %v", err)
  //}

  //conMan, err = util.NewConMan(config)
  //if err != nil {
    //return nil, fmt.Errorf("Cannot create con-man: %v", err)
  //}

  dserver := &Server{}
  
  mux := http.NewServeMux()

  SetupClientApi(dserver, mux)

  server := &http.Server{Handler: mux}

  clientSocketPath := ""

  if config.ClientSocketPath == "" {
    clientSocketPath = "/run/cheshire-chatd.sock"
  }

  os.Remove(clientSocketPath)

  //listener, err = net.Listen("unix", clientSocketPath)
  //if err != nil {
    //return nil, fmt.Errorf("Cannot start client unix-socket listener: %v", err)
  //}

  //defer Listen.Close()

  //log.Info("Starting client api server at %s\n", clientSocketPath)

  dserver.ClientSocketPath = clientSocketPath
  dserver.ClientServer = server

  return dserver, nil
}
