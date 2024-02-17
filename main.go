package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

/*
	TODO: [A] Comment every function
	TODO: [A] take care of every error
	TODO: [A] validate file paths before passing to Encrypt() or Decrypt() (filepath and stat)
	TODO: [A] try to optimize and streamline logging and error output (should be one explicit error format)

	TODO: [B] Try optimizing the code
	TODO: [B] Try implementing wildcards for file names
	TODO: [B] during config check if the files exist. if not > reconfigure
	TODO: [B] implement healthcheck that verifies that everything is working correctly
	TODO: [B] be sure that everything is working correctly (test)
	TODO: [B] make the config stage more user friendly

	TODO: [C] support other key types
	TODO: [C] Set variables to useful names
	TODO: [C] try to make the package smaller when compiled
	TODO: [C] don't hardcode the key type to ssh and ask during config stage
	TODO: [C] maybe change cli framework? I really wanna use something from charm




*/

// Config struct for config File
type Config struct {
	PubKeyPath  string `json:"PubKeyPath"`
	PrivKeyPath string `json:"PrivKeyPath"`
	KeyType     string `json:"KeyType"`
}

// path to config file
const configFilePath = ".agencrc"

// TODO: [B] if no command is passed it should encrypt per default
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
		fmt.Println(err)
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
		fmt.Println("error getting user dir: ", err)
	}

	configPath := filepath.Join(homeDir, configFilePath)

	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("config file exists.")
		return nil
	}

	fmt.Println("enter the path to your public key: ")
	PubkeyPath := ""
	_, err = fmt.Scanln(&PubkeyPath)
	if err != nil {
		return err
	}
	fmt.Println("enter the path to your private key: ")
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
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("config file created.")
	return nil
}
