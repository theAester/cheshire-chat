package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WorkingDir string `yaml:"working_dir"`
	LogFileName string `yaml:"log_file_name"`
  ClientSocketPath string `yaml:"client_socket_path"`
	Text string `yaml:"text"`
}

// parseConfig parses the YAML config file and sets default values for optional fields
func ParseConfig(configFile string) (*Config, error) {
	// Open the config file
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// Read the config file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the YAML content into the Config struct
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set default values for optional fields if they are empty
	setDefaultIfEmpty(&config.WorkingDir, "/opt/cheshire-chatd/default/")
	setDefaultIfEmpty(&config.LogFileName, "")
	setDefaultIfEmpty(&config.Text, "penisman")
	// For slices and nested structs, you would need to check each element individually

	return &config, nil
}


func setDefaultIfEmpty(field *string, defaultValue string) {
	if *field == "" {
		*field = defaultValue
	}
}

