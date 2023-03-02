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
	"github.com/tidwall/gjson"
)

type ProfileInfo struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Picture   string
}

type LinkedinAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

var (
	rootTokenLinkedinURL = "https://www.linkedin.com/oauth/v2/accessToken"
	emailInfoURL         = "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))&oauth2_access_token="
	userInfoURL          = "https://api.linkedin.com/v2/me"
	userPicURL           = "https://api.linkedin.com/v2/me?projection=(id,firstName,lastName,profilePicture(displayImage~:playableStreams))"
)

func GetLinkedinOauthToken(code string, cfg config.Config) (*LinkedinAccessToken, error) {
	values := url.Values{}

	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", cfg.Auth.LinkedinAuthClientID)
	values.Add("client_secret", cfg.Auth.LinkedinAuthClientSecret)
	values.Add("redirect_uri", cfg.Auth.LinkedinAuthRedirectURL)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootTokenLinkedinURL, bytes.NewBufferString(query))

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

	var TokenRes map[string]interface{}

	err = json.Unmarshal(resBody.Bytes(), &TokenRes)
	if err != nil {
		return nil, err
	}

	token := &LinkedinAccessToken{
		AccessToken: TokenRes["access_token"].(string),
		Scope:       TokenRes["scope"].(string),
	}

	return token, nil
}

func GetLinkedinUser(cfg config.Config, token *LinkedinAccessToken) (*ProfileInfo, error) {
	var userProfileInfo ProfileInfo

	// getting user email
	reqUserEmail, err := http.NewRequest("GET", emailInfoURL, nil)

	if err != nil {
		return nil, err
	}

	authorizationHeaderValue := fmt.Sprintf("Bearer %s", token.AccessToken)
	reqUserEmail.Header.Set("Authorization", authorizationHeaderValue)

	clientUserEmail := http.Client{
		Timeout: time.Second * 30,
	}

	resUserEmail, err := clientUserEmail.Do(reqUserEmail)
	if err != nil {
		return nil, err
	}

	defer resUserEmail.Body.Close()

	if resUserEmail.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user email")
	}

	contentUserEmail, err := io.ReadAll(resUserEmail.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.Email = gjson.Get(string(contentUserEmail), "elements.0.handle~.emailAddress").String()

	// getting user main info
	reqUserInfo, err := http.NewRequest("GET", userInfoURL, nil)

	if err != nil {
		return nil, err
	}

	reqUserInfo.Header.Set("Authorization", authorizationHeaderValue)

	clientUserInfo := http.Client{
		Timeout: time.Second * 30,
	}

	resUserInfo, err := clientUserInfo.Do(reqUserInfo)
	if err != nil {
		return nil, err
	}

	defer resUserInfo.Body.Close()

	if resUserInfo.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user info")
	}

	contentUserInfo, err := io.ReadAll(resUserInfo.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.ID = gjson.Get(string(contentUserInfo), "id").String()
	userProfileInfo.FirstName = gjson.Get(string(contentUserInfo), "localizedFirstName").String()
	userProfileInfo.LastName = gjson.Get(string(contentUserInfo), "localizedLastName").String()

	//  get user pic
	reqUserPic, err := http.NewRequest("GET", userPicURL, nil)

	if err != nil {
		return nil, err
	}

	reqUserPic.Header.Set("Authorization", authorizationHeaderValue)

	clientUserPic := http.Client{
		Timeout: time.Second * 30,
	}

	resUserPic, err := clientUserPic.Do(reqUserPic)
	if err != nil {
		return nil, err
	}

	defer resUserPic.Body.Close()

	if resUserPic.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user picture")
	}

	contentUserPic, err := io.ReadAll(resUserPic.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.Picture = gjson.Get(string(contentUserPic), "profilePicture.displayImage~.elements.0.identifiers.0.identifier").String()

	return &userProfileInfo, nil
}
