package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

/*
	TODO: [ ] Comment every function
	TODO: [ ] take care of every error
	TODO: [ ] validate file paths before passing to Encrypt() or Decrypt()
	TODO: [ ] Validate if the files are actually encrypted before Decryption (using the age header)
	TODO: [ ] Validate if the files are actually unencrypted before Encryption (don't know how to do that yet)
	TODO: [ ] if no command is passed it should encrypt per default
	TODO: [ ] Try optimizing the code
	TODO: [ ] try to optimize and streamline logging and error output
	TODO: [ ] try to make the package smaller when compiled
	TODO: [ ] Try implementing wildcards
	TODO: [ ] be sure that everything is working correctly
	TODO: [ ] support other key types
	TODO: [ ] don't hardcode the key type to ssh and ask during config stage
	TODO: [ ] make the config stage more user friendly
	TODO: [ ] during config check if the files exist. if not > reconfigure
	TODO: [ ] implement healthcheck that verifies that everything is working correctly
	TODO: [ ] maybe change cli framework? I really wanna use something from charm




*/

// Config struct for config File
type Config struct {
	PubKeyPath  string `json:"PubKeyPath"`
	PrivKeyPath string `json:"PrivKeyPath"`
	KeyType     string `json:"KeyType"`
}

// path to config file
const configFilePath = ".agencrc"

func main() {
	app := &cli.App{
		Name:  "agenc",
		Usage: "encrypt and decrypt using a preset key in age for convenience",
		Commands: []*cli.Command{
			{
				Name:    "encrypt",
				Aliases: []string{"enc", "e"},
				Usage:   "encrypt",
				Action:  encryptAction,
			},
			{
				Name:    "decrypt",
				Aliases: []string{"dec", "d"},
				Usage:   "decrypt",
				Action:  decryptAction,
			},
			{
				Name:    "config",
				Aliases: []string{"conf", "c"},
				Usage:   "Set config file",
				Action:  configAction,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func encryptAction(c *cli.Context) error {
	config, _ := readConfig()
	recipient := getRecipientFromSSHKeyFile(config)
	for i := 0; i < c.Args().Len(); i++ {
		encryptFile(c.Args().Get(i), recipient)
	}
	return nil
}

func decryptAction(c *cli.Context) error {
	config, _ := readConfig()
	identity, _ := getIdentityFromSSHKeyFile(config)
	for i := 0; i < c.Args().Len(); i++ {
		decryptFile(c.Args().Get(i), identity)
	}
	return nil
}

func configAction(c *cli.Context) error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(homeDir, configFilePath)

	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("config file exists.")
		return nil
	}

	fmt.Print("enter the path to your public key: ")
	PubkeyPath := ""
	_, err = fmt.Scanln(&PubkeyPath)
	if err != nil {
		return err
	}
	fmt.Print("enter the path to your private key: ")
	PrivkeyPath := ""
	_, err = fmt.Scanln(&PrivkeyPath)
	if err != nil {
		return err
	}
	config := Config{
		PubKeyPath:  PubkeyPath,
		PrivKeyPath: PrivkeyPath,
		KeyType:     "SSH",
	}

	file, err := os.Create(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("config file created.")
	return nil
}
