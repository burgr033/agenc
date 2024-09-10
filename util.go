package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var expectedHeaderHex = "6167652D656E6372797074696F6E2E6F72672F76310A" // age-encryption.org/v1

// getOutputFileName converts the input filename to the output filename (adds the age extension)
func getOutputFileName(inputFilename string) string {
	outputFilename := inputFilename + ".age"
	return outputFilename
}

// checkAgeHeaderPresence takes a file name opens the file and checks the first 22 bytes for the age header
func checkAgeHeaderPresence(inputFilename string) bool {
	header, _ := os.Open(inputFilename)
	expectedHeader, _ := hex.DecodeString(expectedHeaderHex)
	headerBytes := make([]byte, len(expectedHeader))
	header.Read(headerBytes)
	if string(headerBytes) == string(expectedHeader) {
		header.Close()
		return true
	} else {
		header.Close()
		return false
	}
}

// getOutputFileNameDecrypt converts the taken file name and trims .age if existing
func getOutputFileNameWithoutAge(inputFilename string) string {
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

// getSSHKeyFileContent is given the config and returns the private key as byte slice
func getSSHKeyFileContent(config Config, isPublic bool) ([]byte, error) {
	var filepath string
	if isPublic {
		filepath = config.PubKeyPath
	} else {
		filepath = config.PrivKeyPath
	}

	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return b, nil
}
