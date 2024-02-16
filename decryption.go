package main

import (
	"filippo.io/age"
	"filippo.io/age/agessh"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"io"
	"os"
)

func getIdentityFromSSHKeyFile(config Config) (age.Identity, error) {
	keyContent := getSSHPrivKeyFileContent(config)
	fmt.Print("Please enter SSH Key Password for ", config.PrivKeyPath, ": ")
	pwBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
	clearTextKey, _ := ssh.ParseRawPrivateKeyWithPassphrase(keyContent, pwBytes)
	switch key := clearTextKey.(type) {
	case *ed25519.PrivateKey:
		return agessh.NewEd25519Identity(*key)
	default:
		return nil, fmt.Errorf("ffs")
	}
}
func decryptFile(inputFilename string, identity age.Identity) {
	outputFilename := getOutputFileNameDecrypt(inputFilename)
	encryptedFile, _ := os.Open(inputFilename)
	defer encryptedFile.Close()
	decryptedReader, _ := age.Decrypt(encryptedFile, identity)
	outputFile, _ := os.Create(outputFilename)
	defer outputFile.Close()
	io.Copy(outputFile, decryptedReader)
	encryptedFile.Close()
	outputFile.Close()
}

// reads config file
func getSSHPrivKeyFileContent(config Config) []byte {
	filepath := config.PrivKeyPath
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
