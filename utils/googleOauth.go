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

const OAuthStateCookieName string = "oauthstate"

type GoogleOauthToken struct {
	AccessToken string
	IDToken     string
}

func GetGoogleOuathToken(code string, cfg config.Config) (*GoogleOauthToken, error) {
	const rootURL = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", cfg.Auth.GoogleAuthClientID)
	values.Add("client_secret", cfg.Auth.GoogleAuthClientSecret)
	values.Add("redirect_uri", cfg.Auth.GoogleAuthRedirectURL)

	query := values.Encode()

	req, err := http.NewRequest("POST", rootURL, bytes.NewBufferString(query))
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

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		AccessToken: GoogleOauthTokenRes["access_token"].(string),
		IDToken:     GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

type GoogleUserResult struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Locale        string
	Picture       string
}

func GetGoogleUser(accessToken, idToken string) (*GoogleUserResult, error) {
	rootURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken)

	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		ID:            GoogleUserRes["id"].(string),
		Email:         GoogleUserRes["email"].(string),
		VerifiedEmail: GoogleUserRes["verified_email"].(bool),
		Name:          GoogleUserRes["name"].(string),
		Picture:       GoogleUserRes["picture"].(string),
		GivenName:     GoogleUserRes["given_name"].(string),
		Locale:        GoogleUserRes["locale"].(string),
	}

	return userBody, nil
}

//func GenerateStateOauthCookie(w http.ResponseWriter, maxAge int, path, domain string) string {
//	b := make([]byte, 16)
//	rand.Read(b)
//	state := base64.URLEncoding.EncodeToString(b)
//	cookie := http.Cookie{}
//	cookie.Name = OAuthStateCookieName
//	cookie.Value = state
//	cookie.Path = path
//	cookie.Domain = domain
//	cookie.MaxAge = maxAge
//	cookie.Secure = false
//	cookie.HttpOnly = true
//	http.SetCookie(w, &cookie)
//
//	return state
//}
