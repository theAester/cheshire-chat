////////// CHESHIRE CHAT DAEMON //////////
// In this file we initialize and start //
// The daemon.				//
/////////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/theAester/cheshire-chat/cmd/daemon/util"
	"github.com/theAester/cheshire-chat/cmd/daemon/server"
)

// initLog initializes the log module to write to the specified log file or stdout if not specified
func initLog(config *util.Config) {
	if config.LogFileName != "" {
		logFile, err := os.OpenFile(filepath.Join(config.WorkingDir, config.LogFileName), 
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open %s for writing logs: %v\n", config.LogFileName, err);
			os.Exit(1)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
}

func main(){
	// Define flag for config file path
	configFile := *flag.String("c", "/etc/cheshire-chatd/config.yaml", "Path to config file")
	flag.Parse()

	daemonConfig, err := util.ParseConfig(configFile)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error while parsing config file: %v\n", err)
		os.Exit(1)
	}

	initLog(daemonConfig)

	_, err = server.New(daemonConfig)
	if err != nil {
		log.Fatalf("Error while starting server: %v\n", err)
	}
}
