package yauth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Token holds all values nececcary to access
// Yandex services using Yandex OAuth.
type Token struct {
	// Access field contains the actual access token
	Access string `json:"access_token,omitempty"`
	// Refresh field containes refresh token that
	// can be used to update access token
	Refresh string `json:"refresh_token,omitempty"`
	// Expiry holds token expiration time
	Expiry time.Time `json:"token_expires,omitempty"`
}

// OpenToken tries to read filename and parse token values.
func OpenToken(filename string) (Token, error) {
	var t Token
	data, err := os.ReadFile(filename)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(data, &t)
	return t, err
}

// SaveToFile marshals token to json and writes to filename.
func (t *Token) SaveToFile(filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return err
	}
	data, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Token) Valid() bool {
	if t == nil || t.Access == "" || t.Expiry.Before(time.Now()) {
		return false
	}
	return true
}

func (t *Token) String() string {
	return fmt.Sprintf("%s\nexpires on: %s", t.Access, t.Expiry.Format("15:04:05 02 Jan 2006"))
}
