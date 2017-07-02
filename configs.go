package gfs

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"log"
	"os"
)

// A config value
type Config struct {
	// The username required when uploading files
	Username string `json:"username"`
	// The password required when uploading files
	Password string `json:"password"`
	// The path that should be served
	Serve string `json:"serve"`
	// The port to serve on
	Port string `json:"port"`
	// The secret used to verify authorized requests
	Secret string `json:"secret"`
	// Indicates if login is required to be allowed to read the contents
	LoginRequiredForRead bool `json:"loginRequiredForRead"`
}

// Reads the specified config file
func readConfigFile(path string) (config *Config, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config = new(Config)
	err = json.NewDecoder(file).Decode(config)
	return config, err
}

// Reads the config files, or gives default values if the config
// doesn't yet exist
func GetConfigs(path string) (config *Config, err error) {
	config, err = readConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Config file didn't exist, defaulting")

			password, err := CreatePassword("password")
			if err != nil {
				return nil, err
			}

			config = &Config{
				Username: "username",
				Password: password,
				Serve:    DefaultServePath,
				Port:     "8080",
				Secret:   uuid.NewV4().String(),
			}

			SaveConfigs(path, config)

			return config, nil
		}
		return nil, err
	}

	return config, nil
}

// Saves the given config to disk
func SaveConfigs(path string, config *Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(*config)
}
