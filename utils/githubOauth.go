package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/red-rocket-software/reminder-go/config"
)

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUserResult struct {
	Name    string
	Email   string
	Picture string
}

const (
	rootTokenURL = "https://github.com/login/oauth/access_token"
	rootUserURL  = "https://api.github.com/user"
)

func GetGithubOuathToken(code string, cfg config.Config) (*GithubAccessToken, error) {
	values := url.Values{}

	values.Add("code", code)
	values.Add("client_id", cfg.Auth.GithubAuthClientID)
	values.Add("client_secret", cfg.Auth.GithubAuthClientSecret)
	values.Add("redirect_uri", cfg.Auth.GithubAuthRedirectURL)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootTokenURL, bytes.NewBufferString(query))

	req, err := http.NewRequest("POST", queryString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)

	if err != nil {
		return nil, err
	}

	parsedQuery, err := url.ParseQuery(resBody.String())

	if err != nil {
		return nil, err
	}

	token := &GithubAccessToken{
		AccessToken: parsedQuery["access_token"][0],
		TokenType:   parsedQuery["token_type"][0],
		Scope:       parsedQuery["scope"][0],
	}

	return token, nil
}

func GetGithubUser(token *GithubAccessToken) (*GithubUserResult, error) {
	req, err := http.NewRequest("GET", rootUserURL, nil)

	if err != nil {
		return nil, err
	}

	authorizationHeaderValue := fmt.Sprintf("Bearer %s", token.AccessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, errors.New("could not parse response")
	}

	var GithubUserRes map[string]interface{}

	err = json.Unmarshal(resBody.Bytes(), &GithubUserRes)
	if err != nil {
		return nil, err
	}

	userBody := &GithubUserResult{
		Email:   GithubUserRes["email"].(string),
		Name:    GithubUserRes["name"].(string),
		Picture: GithubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}
