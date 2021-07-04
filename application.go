package yauth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Application represents app registered with Yandex OAuth
type Application struct {
	ClientID string `json:"client_id,omitempty"`
	Secret   string `json:"client_secret,omitempty"`
	filename string
}

// New creates Application with given credentials
func New(clientID, secret string) (*Application, error) {
	var app Application
	var err error
	if clientID == "" || secret == "" {
		err = fmt.Errorf("error creating application: id and password can not be empty")
		return &app, err
	}
	app.ClientID = clientID
	app.Secret = secret
	return &app, err
}

// OpenApp accepts name of file storing application credentials
// in JSON format
func Open(filename string) (*Application, error) {
	app := &Application{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return app, err
	}
	err = json.Unmarshal(data, app)
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
// configured scope.
func (app *Application) RequestUserAuthorization() (accesstoken, refreshtoken string, expires int, err error) {
	accesstoken, refreshtoken, expires, err = requestUserAuthorization(app.ClientID, app.Secret)
	return
}

// RefreshToken accepts refresh token, connects to Yandex OAuth API,
// requests new access token and returns it.
func (app *Application) RefreshToken(refresh string) (accesstoken, refreshtoken string, expires int, err error) {
	accesstoken, refreshtoken, expires, err = renewToken(app.ClientID, app.Secret, refresh)
	return
}

// SaveToFile saves app credentials to specified filename. It is intended
// to be used after creating app with yauth.New(). If this method returned no
// error, the App can later be saved with Save() method.
func (app *Application) SaveToFile(filename string) error {
	err := saveApptoFile(app, filename)
	if err == nil {
		app.filename = filename
	}
	return err
}

// Save saves application to the same file from which it was read with yauth.Open()
// function.
func (app *Application) Save() error {
	return saveApptoFile(app, app.filename)
}

func saveApptoFile(app *Application, filename string) error {
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, 0755)
	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
