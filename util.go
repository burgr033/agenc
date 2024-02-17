package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// converts the input filename to the output filename (adds the age extension)
func getOutputFileName(inputFilename string) string {

	outputFilename := inputFilename + ".age"
	return outputFilename

}

// TODO: [A] Currently the file suffix is just trimmed anyway, without checking that it's actually the age suffix. It should be checked if there is an age suffix. If not this function should return the normal file name.
// TODO: [B] better function name
func getOutputFileNameDecrypt(inputFilename string) string {

	extension := filepath.Ext(inputFilename)
	outputFilename := strings.TrimSuffix(inputFilename, extension)
	return outputFilename

}
func readConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configPath := filepath.Join(homeDir, configFilePath)

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
