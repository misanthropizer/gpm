package main

//////////////////////////////////////////////////////////////////
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// ////////////////////////////////////////////////////////////////
type Config struct {
	PrivateKey      string `json:"private_key"`
	PublicKey       string `json:"public_key"`
	PastebinAPIKey  string `json:"pastebin_api_key"`
	PastebinUserKey string `json:"pastebin_user_key"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}

const configFilePath = "./config.json"

// ///////////////////////////CONFIG//////////////////////////////
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

// /////////////////////////GPG HANDLING////////////////////
func encryptMessage(publicKey string, message string) (string, error) {
	return helper.EncryptMessageArmored(publicKey, message)
}

func decryptMessage(privateKey string, passphrase string, encryptedMessage string) (string, error) {
	return helper.DecryptMessageArmored(privateKey, []byte(passphrase), encryptedMessage)
}

// ///////////////////DATA EXTRACTION/////////////////////////////////////////////
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

// ///////////////////////////////PASTEBIN/////////////////////////////////////////////////////////
func postToPastebin(apiKey, userKey, encryptedMessage string) (string, error) {
	data := url.Values{
		"api_dev_key":       {apiKey},
		"api_user_key":      {userKey},
		"api_option":        {"paste"},
		"api_paste_code":    {encryptedMessage},
		"api_paste_private": {"1"}, // 1 = unlisted, 2 = private
	}

	resp, err := http.PostForm("https://pastebin.com/api/api_post.php", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func getPastebinUserKey(apiKey, username, password string) (string, error) {
	data := url.Values{
		"api_dev_key":       {apiKey},
		"api_user_name":     {username},
		"api_user_password": {password},
	}

	resp, err := http.PostForm("https://pastebin.com/api/api_login.php", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// /////////////////////////////MAIN///////////////////////////////////
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

		fmt.Println("Enter your Pastebin API key: ")
		config.PastebinAPIKey, _ = reader.ReadString('\n')
		config.PastebinAPIKey = strings.TrimSpace(config.PastebinAPIKey)
		fmt.Println("Enter your Pastebin username: ")
		config.Username, _ = reader.ReadString('\n')
		config.Username = strings.TrimSpace(config.Username)
		fmt.Println("Enter your Pastebin password: ")
		config.Password, _ = reader.ReadString('\n')
		config.Password = strings.TrimSpace(config.Password)
		if err := saveConfig(&config); err != nil {
			fmt.Println("Failed to save config:", err)
			return
		}
		fmt.Println("Config saved.")
	}
	/////////////////ERROR/////////////////////////////////////////////////
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
	config.PastebinUserKey, err = getPastebinUserKey(config.PastebinAPIKey, config.Username, config.Password)
	if err != nil {
		fmt.Println("Failed to get Pastebin user key:", err)
		return
	}
	/////////////////CLI MENU/////////////////////////////////////////////////
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

			pastebinLink, err := postToPastebin(config.PastebinAPIKey, config.PastebinUserKey, encryptedMessage)
			if err != nil {
				fmt.Println("Failed to post to Pastebin:", err)
				continue
			}
			fmt.Println("Message posted to Pastebin:", pastebinLink)
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
			fmt.Println("Invalid option")
		}
	}
}
