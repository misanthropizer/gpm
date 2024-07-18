package main

//////////////////////////////////////////////////////////////////
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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

func getUserPastes(apiKey, userKey string) ([]string, error) {
	data := url.Values{
		"api_dev_key":       {apiKey},
		"api_user_key":      {userKey},
		"api_option":        {"list"},
		"api_results_limit": {"100"},
	}

	resp, err := http.PostForm("https://pastebin.com/api/api_post.php", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pasteIDs []string
	pasteEntries := strings.Split(string(body), "<paste_key>")
	for _, entry := range pasteEntries[1:] {
		keyEnd := strings.Index(entry, "</paste_key>")
		if keyEnd == -1 {
			continue
		}
		pasteIDs = append(pasteIDs, entry[:keyEnd])
	}

	return pasteIDs, nil
}

func getPasteContent(pasteID string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://pastebin.com/raw/%s", pasteID))
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
