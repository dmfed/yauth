package yauth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Application represents app registered with Yandex OAuth
type Application struct {
	ClientID string `json:"client_id"`
	Secret   string `json:"client_secret"`
	filename string
}

// NewApp creates Application with given credentials
func NewApp(clientID, secret string) (Application, error) {
	var app Application
	var err error
	if clientID == "" || secret == "" {
		err = fmt.Errorf("error creating application: id and password can not be empty")
		return app, err
	}
	app.ClientID = clientID
	app.Secret = secret
	return app, err
}

// OpenApp accepts name of file storing application credentials
// in JSON format:
// {
//    "client_id": "your_app_id",
//    "client_secret": "your_app_secret"
// }
// If credentials were parsed successfully returns Application.
func OpenApp(filename string) (Application, error) {
	var app Application
	data, err := os.ReadFile(filename)
	if err != nil {
		return app, err
	}
	err = json.Unmarshal(data, &app)
	if err == nil {
		app.filename = filename
	}
	return app, err
}

func (app *Application) String() string {
	data, _ := json.MarshalIndent(app, "", "    ")
	return string(data)
}

// RequestUserAuthorization connects to Yandex OAuth API, fetches
// device code and user code and then asks user to follow a link on
// any of their devices to authorize application to access its default
// scope configured with Yandex OAuth.
// It outputs the URL to follow and code to os.Stdout.
func (app *Application) RequestUserAuthorization() (token Token, err error) {
	accesstoken, refreshtoken, expires, err := requestUserAuthorization(app.ClientID, app.Secret)
	if err == nil {
		token.Access = accesstoken
		token.Refresh = refreshtoken
		token.Expiry = time.Now().Add(time.Second * time.Duration(expires))
	}
	return
}

// Refresh connects to Yandex OAuth API,
// requests new access token and saves it to Application.
func (app *Application) Refresh(refreshtoken string) (token Token, err error) {
	if refreshtoken == "" {
		return token, fmt.Errorf("refreshtoken can not be empty")
	}
	access, refresh, expires, err := renewToken(app.ClientID, app.Secret, refreshtoken)
	if err == nil {
		token.Access = access
		token.Refresh = refresh
		token.Expiry = time.Now().Add(time.Second * time.Duration(expires))
	}
	return
}

// SaveToFile saves app credentials to specified filename. It is intended
// to be used after creating app with yauth.New(). If this method returned no
// error, the App can later be saved with Save() method.
func (app *Application) SaveToFile(filename string) error {
	return saveApptoFile(app, filename)
}

// Save saves Application state to disk. If will return error if
// Application was just created (not opened with Open). Use SaveToFile method
// to save your app for the first time.
func (app *Application) Save() error {
	if app.filename == "" {
		return fmt.Errorf("app filename unknown. use SaveToFile method")
	}
	return saveApptoFile(app, app.filename)
}

func saveApptoFile(app *Application, filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return err
	}
	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
