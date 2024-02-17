package main

import (
	"encoding/hex"
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
	keyContent := getSSHPrivKeyFileContent(config)
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
// TODO: [A] Clean Up this function and put it the header stuff into separate function so that we can use it for encryption aswell
func decryptFile(inputFilename string, identity age.Identity) {
	outputFilename := getOutputFileNameDecrypt(inputFilename)
	encryptedFile, _ := os.Open(inputFilename)
	//header stuff
	header, _ := os.Open(inputFilename)
	expectedHeaderHex := "6167652D656E6372797074696F6E2E6F72672F76310A2D3E207373682D65643235353139"
	expectedHeader, _ := hex.DecodeString(expectedHeaderHex)
	headerBytes := make([]byte, len(expectedHeader))
	header.Read(headerBytes)
	if string(headerBytes) == string(expectedHeader) {
		fmt.Println("headers match, so file is encrypted")
	} else {
		fmt.Println("File is not encrypted (no header present)")
		os.Exit(1)
	}
	header.Close()
	//header stuff end
	defer encryptedFile.Close()
	decryptedReader, _ := age.Decrypt(encryptedFile, identity)
	outputFile, _ := os.Create(outputFilename)
	defer outputFile.Close()
	io.Copy(outputFile, decryptedReader)
	encryptedFile.Close()
	outputFile.Close()
}

// TODO: [B] Check if we could merge this with getSSHPubKeyFileContent which returns a string and not byte...
// getSSHPrivKeyFileContent is given the config and returns the private key as byte slice
func getSSHPrivKeyFileContent(config Config) []byte {
	filepath := config.PrivKeyPath
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
