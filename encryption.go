package main

import (
	"bytes"
	"filippo.io/age"
	"filippo.io/age/agessh"
	"fmt"
	"io"
	"os"
)

// uses ParseRecipient from agessh to return recipient
func getRecipientFromSSHKeyFile(config Config) age.Recipient {
	keyContent := getSSHPubKeyFileContent(config)
	recipient, _ := agessh.ParseRecipient(keyContent)
	return recipient
}

// reads Pubkey file
func getSSHPubKeyFileContent(config Config) string {
	filepath := config.PubKeyPath
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

// encryptFile encrypts the file via age and saves the output to the new file
func encryptFile(inputFilename string, recipient age.Recipient) {

	var inputFile, _ = os.Open(inputFilename)
	var inputData bytes.Buffer
	var outputFilename = getOutputFileName(inputFilename)
	var outputFile, _ = os.Create(outputFilename)

	defer outputFile.Close()

	io.Copy(&inputData, inputFile)

	w, _ := age.Encrypt(outputFile, recipient)

	io.Copy(w, &inputData)

	w.Close()
}
