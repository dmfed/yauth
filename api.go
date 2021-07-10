package yauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	baseURL       string = "https://oauth.yandex.ru"
	deviceCodeURL        = baseURL + "/device/code"
	tokenURL             = baseURL + "/token"
)

type CodesResponse struct {
	DeviceCode      string `json:"device_code,omitempty"`
	UserCode        string `json:"user_code,omitempty"`
	VerificationURL string `json:"verification_url,omitempty"`
	Interval        int    `json:"interval,omitempty"`
	ExpiresIn       int    `json:"expires_in,omitempty"`
}

type OAuthAPIError struct {
	Desc string `json:"error_description,omitempty"`
	Err  string `json:"error,omitempty"`
}

func (e OAuthAPIError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.Desc)
}

type Token struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// requestUserAuthorization accepts application ID and password, requests
// authorization codes for device then tells user in current terminal session
// to follow authorization link in a browser and enter device code.
// deviceName is passed to Yandex oauth service if not empty.
// The function returns access token, refresh token, token expiry time in seconds
// and error if something went wrong
func requestUserAuthorization(clientID, clientSecret string) (accesstoken, refreshtoken string, expires int, err error) {
	d := newDevice()
	codes, err := fetchAuthorizationCodes(clientID, d.deviceName)
	if err != nil {
		return
	}
	askUserToFollowURL(codes.VerificationURL, codes.UserCode)
	token, err := fetchOAuthToken(clientID, clientSecret, codes.DeviceCode, codes.Interval+1, codes.ExpiresIn)
	if err == nil {
		accesstoken = token.AccessToken
		refreshtoken = token.RefreshToken
		expires = token.ExpiresIn
	}
	return
}

// renewToken accepts application ID, password and refresh token and requests
// new access token and refresh token from Yandex oauth service.
func renewToken(clientID, clientSecret, refreshToken string) (accesstoken, refreshtoken string, expires int, err error) {
	if refreshToken == "" || clientID == "" || clientSecret == "" {
		err = fmt.Errorf("credentials missing or no refresh token present")
		return
	}
	token, err := requestRefreshToken(clientID, clientSecret, refreshToken)
	if err == nil {
		accesstoken = token.AccessToken
		refreshtoken = token.RefreshToken
		expires = token.ExpiresIn
	}
	return
}

// fetchAuthorizationCodes connects to YAndex OAuth endpoint and requests authorization codes
func fetchAuthorizationCodes(clientID, deviceName string) (codes CodesResponse, err error) {
	v := url.Values{}
	v.Set("client_id", clientID)
	if deviceName != "" {
		v.Set("device_name", deviceName)
	}
	resp, err := http.PostForm(deviceCodeURL, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = extractOAuthAPIError(b)
		return
	}
	err = json.Unmarshal(b, &codes)
	return
}

// fetchOAuthToken polls Yandex OAuth server with interval (in seconds) until device code
// expiry. It returns Token if successful.
func fetchOAuthToken(clientID, clientSecret, deviceCode string, interval, expires int) (t Token, err error) {
	expiry := time.NewTimer(time.Duration(expires) * time.Second)
	retry := time.NewTimer(time.Duration(interval) * time.Second)
	for {
		select {
		case <-expiry.C:
			return
		case <-retry.C:
			retry.Reset(time.Duration(interval) * time.Second)
			if t, err = requestToken(clientID, clientSecret, deviceCode); err == nil {
				return
			}
		}
	}
}

func requestToken(clientID, clientSecret, deviceCode string) (t Token, err error) {
	v := url.Values{}
	v.Set("grant_type", "device_code")
	v.Set("code", deviceCode)
	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	resp, err := http.PostForm(tokenURL, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = extractOAuthAPIError(b)
		return
	}
	if err = json.Unmarshal(b, &t); err != nil {
		return
	}
	return
}

func requestRefreshToken(clientID, clientSecret, refreshToken string) (t Token, err error) {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)
	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	resp, err := http.PostForm(tokenURL, v)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = extractOAuthAPIError(b)
		return
	}
	err = json.Unmarshal(b, &t)
	return
}

func extractOAuthAPIError(b []byte) error {
	e := OAuthAPIError{}
	err := json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	return e
}

func askUserToFollowURL(url, code string) {
	fmt.Printf("Please open this link in a browser and follow the instructions to authorize access to your data:\n%v\n", url)
	fmt.Println("Enter the following verification code when prompted:", code)
}
