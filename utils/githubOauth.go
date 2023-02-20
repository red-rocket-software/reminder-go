package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/red-rocket-software/reminder-go/config"
)

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUserResult struct {
	ID      string
	Name    string
	Email   string
	Picture string
}

const (
	rootTokenURL = "https://github.com/login/oauth/access_token"
	rootUserURL  = "https://api.github.com/user"
)

func GetGithubOuathToken(code string, cfg config.Config) (*GithubAccessToken, error) {
	clientID := cfg.Auth.GithubAuthClientID

	clientSecret := cfg.Auth.GithubAuthClientSecret

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}

	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest("POST", rootTokenURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var token GithubAccessToken

	json.Unmarshal(respBody, &token)

	return &token, nil

}

func GetGithubUser(token *GithubAccessToken) (*GithubUserResult, error) {
	req, err := http.NewRequest("GET", rootUserURL, nil)

	if err != nil {
		return nil, err
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", token.AccessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, errors.New("could not parse response")
	}

	var GithubUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GithubUserRes); err != nil {
		return nil, err
	}

	userBody := &GithubUserResult{
		ID:      GithubUserRes["id"].(string),
		Email:   GithubUserRes["email"].(string),
		Name:    GithubUserRes["name"].(string),
		Picture: GithubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}
