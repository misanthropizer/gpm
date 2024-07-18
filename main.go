package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

type Config struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

const configFilePath = "./config.json"

func loadConfig() (*Config, error) {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath, data, 0644)
}

func encryptMessage(publicKey string, message string) (string, error) {
	return helper.EncryptMessageArmored(publicKey, message)
}

func decryptMessage(privateKey string, passphrase string, encryptedMessage string) (string, error) {
	return helper.DecryptMessageArmored(privateKey, []byte(passphrase), encryptedMessage)
}

func getNameFromKey(key string) (string, error) {
	keyObj, err := crypto.NewKeyFromArmored(key)
	if err != nil {
		return "", err
	}
	entity := keyObj.GetEntity()
	if len(entity.Identities) > 0 {
		for _, id := range entity.Identities {
			return id.Name, nil
		}
	}
	return "", fmt.Errorf("no identity found in key")
}

func main() {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Config file not found. Creating a new one.")

		var config Config
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Enter your private key (end with a single dot '.'): ")
		var privateKey strings.Builder
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "." {
				break
			}
			privateKey.WriteString(line + "\n")
		}
		config.PrivateKey = privateKey.String()

		fmt.Println("Enter your public key (end with a single dot '.'): ")
		var publicKey strings.Builder
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "." {
				break
			}
			publicKey.WriteString(line + "\n")
		}
		config.PublicKey = publicKey.String()

		if err := saveConfig(&config); err != nil {
			fmt.Println("Failed to save config:", err)
			return
		}
		fmt.Println("Config saved.")
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	privateKeyName, err := getNameFromKey(config.PrivateKey)
	if err != nil {
		fmt.Println("Failed to get name from private key:", err)
		return
	}
	publicKeyName, err := getNameFromKey(config.PublicKey)
	if err != nil {
		fmt.Println("Failed to get name from public key:", err)
		return
	}

	fmt.Println("Welcome to GPG Chat")
	fmt.Println("You:", privateKeyName)
	fmt.Println("Chat Partner:", publicKeyName)

	for {
		fmt.Println("1 | Send Message")
		fmt.Println("2 | Receive Message")
		fmt.Print("Enter an option: ")

		var option int
		fmt.Scan(&option)

		switch option {
		case 1:
			fmt.Print("Enter message to send: ")
			var message string
			reader := bufio.NewReader(os.Stdin)
			message, _ = reader.ReadString('\n')
			encryptedMessage, err := encryptMessage(config.PublicKey, strings.TrimSpace(message))
			if err != nil {
				fmt.Println("Failed to encrypt message:", err)
				continue
			}
			fmt.Println("Encrypted message:", encryptedMessage)
		case 2:
			fmt.Print("Enter passphrase for private key: ")
			var passphrase string
			fmt.Scan(&passphrase)
			fmt.Print("Enter message to decrypt: ")
			var encryptedMessage string
			reader := bufio.NewReader(os.Stdin)
			encryptedMessage, _ = reader.ReadString('\n')
			decryptedMessage, err := decryptMessage(config.PrivateKey, passphrase, strings.TrimSpace(encryptedMessage))
			if err != nil {
				fmt.Println("Failed to decrypt message:", err)
				continue
			}
			fmt.Printf("%s: %s\n", publicKeyName, decryptedMessage)
		default:
			fmt.Println("Invalid option, stop fucking with it!")
		}
	}
}
