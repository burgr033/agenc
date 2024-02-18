package main

import (
	"bytes"
	"filippo.io/age"
	"filippo.io/age/agessh"
	"io"
	"os"
)

// getRecipientFromSSHKeyFile is given the config and uses ParseRecipient from agessh to return recipient
func getRecipientFromSSHKeyFile(config Config) age.Recipient {
	keyContent, _ := getSSHKeyFileContent(config, true)
	recipient, _ := agessh.ParseRecipient(string(keyContent))
	return recipient
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
