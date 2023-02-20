package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

type ProfileInfo struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Picture   string
}

var (
	emailInfoURL = "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))&oauth2_access_token="
	userInfoURL  = "https://api.linkedin.com/v2/me"
	userPicURL   = "https://api.linkedin.com/v2/me?projection=(id,firstName,lastName,profilePicture(displayImage~:playableStreams))"
)

func GetLinkedInConfig(cfg config.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  cfg.Auth.LinkedinAuthRedirectURL,
		ClientID:     cfg.Auth.LinkedinAuthClientID,
		ClientSecret: cfg.Auth.LinkedinAuthClientSecret,
		Scopes:       []string{"r_emailaddress", "r_liteprofile"},
		Endpoint:     linkedin.Endpoint,
	}
}

func GetLinkedinOauthToken(code string, cfg config.Config) (*oauth2.Token, error) {
	token, err := GetLinkedInConfig(cfg).Exchange(context.Background(), code)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("wrong code: %s", code))
	}
	return token, nil
}

func GetLinkedinHTTPClient(token *oauth2.Token, cfg config.Config) *http.Client {
	return GetLinkedInConfig(cfg).Client(context.Background(), token)
}

func GetLinkedinUser(cfg config.Config, token *oauth2.Token) (*ProfileInfo, error) {
	var userProfileInfo ProfileInfo

	client := GetLinkedInConfig(cfg).Client(context.Background(), token)

	// get user email
	reqUserEmail, err := http.NewRequest("GET", emailInfoURL, nil)

	if err != nil {
		return nil, err
	}

	reqUserEmail.Header.Set("Bearer", token.AccessToken)

	resUserEmail, err := client.Do(reqUserEmail)
	if err != nil {
		return nil, err
	}

	defer resUserEmail.Body.Close()

	contentUserEmail, err := io.ReadAll(resUserEmail.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.Email = gjson.Get(string(contentUserEmail), "elements.0.handle~.emailAddress").String()

	// get user info
	reqUserInfo, err := http.NewRequest("GET", userInfoURL, nil)

	if err != nil {
		return nil, err
	}

	reqUserInfo.Header.Set("Bearer", token.AccessToken)

	resUserInfo, err := client.Do(reqUserInfo)
	if err != nil {
		return nil, err
	}

	defer resUserEmail.Body.Close()

	contentUserInfo, err := io.ReadAll(resUserInfo.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.ID = gjson.Get(string(contentUserInfo), "id").String()
	userProfileInfo.FirstName = gjson.Get(string(contentUserInfo), "localizedFirstName").String()
	userProfileInfo.LastName = gjson.Get(string(contentUserInfo), "localizedLastName").String()

	// get user pic
	reqUserPic, err := http.NewRequest("GET", userPicURL, nil)

	if err != nil {
		return nil, err
	}

	reqUserPic.Header.Set("Bearer", token.AccessToken)

	resUserPic, err := client.Do(reqUserPic)
	if err != nil {
		return nil, err
	}

	defer resUserPic.Body.Close()

	contentUserPic, err := io.ReadAll(resUserPic.Body)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("could not parse response: %s", err.Error()))
	}

	userProfileInfo.Picture = gjson.Get(string(contentUserPic), "profilePicture.displayImage~.elements.#.identifiers.0.identifier").String()

	return &userProfileInfo, nil
}
