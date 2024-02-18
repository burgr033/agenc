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

// getIdentityFromSSHKeyFile takes the json config and returns the age identity by decrypting the encrypted SSH private
// key by password.
func getIdentityFromSSHKeyFile(config Config) (age.Identity, error) {
	keyContent, _ := getSSHKeyFileContent(config, false)
	fmt.Print("Please enter SSH Key Password for ", config.PrivKeyPath, ": ")
	pwBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
	clearTextKey, _ := ssh.ParseRawPrivateKeyWithPassphrase(keyContent, pwBytes)
	fmt.Print("\n")
	switch key := clearTextKey.(type) {
	case *ed25519.PrivateKey:
		return agessh.NewEd25519Identity(*key)
	default:
		return nil, fmt.Errorf("ffs")
	}
}

// decryptFile takes the file name and age identity to decrypt the file and write it to the output file.
// Output File is the file without the .age suffix
func decryptFile(inputFilename string, identity age.Identity) {
	outputFilename := getOutputFileNameWithoutAge(inputFilename)
	encryptedFile, _ := os.Open(inputFilename)
	defer encryptedFile.Close()
	decryptedReader, _ := age.Decrypt(encryptedFile, identity)
	outputFile, _ := os.Create(outputFilename)
	defer outputFile.Close()
	io.Copy(outputFile, decryptedReader)
	encryptedFile.Close()
	outputFile.Close()
}
