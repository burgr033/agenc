package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getOutputFileName converts the input filename to the output filename (adds the age extension)
func getOutputFileName(inputFilename string) string {

	outputFilename := inputFilename + ".age"
	return outputFilename

}

// getOutputFileNameDecrypt converts the taken file name and trims .age if existing
// TODO: [B] better function name
func getOutputFileNameDecrypt(inputFilename string) string {

	extension := filepath.Ext(inputFilename)
	if extension != ".age" {
		fmt.Println("Note: File name does not contain .age extension. Filename not changed.")
		return inputFilename
	}
	outputFilename := strings.TrimSuffix(inputFilename, extension)
	return outputFilename

}

// readConfig reads the config file and returns config type
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
