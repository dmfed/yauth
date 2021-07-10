# Package yauth
import "github.com/dmfed/yauth"

## Usage

First you need to register your application at https://oauth.yandex.ru

Use the credentials received from Yandex oauth service to create yauth.Application as follows:

```go
app, err := yauth.NewApp("your_client_id", "your_secret")
```

Application can be saved to disk and reused later:

```go 
app.SaveToFile("./myapp.json")

app, err := yauth.OpenApp("./myapp.json")
```
Application then can be used to request user authorization as follows:

```go
token, err := app.RequestUserAuthorization()
```
The app then prints instructions to user to Stdout and starts polling 
Yandex OAuth API to check whether user has authorized your application. 

If authorization was successful the app returns AuthToken.

## AuthToken

The AuthToken is defined as follows:

```go
// AuthToken holds all values nececcary to access
// Yandex services using Yandex OAuth.
type AuthToken struct {
	// Access field contains the actual access token
	Access string `json:"access_token,omitempty"`
	// Refresh field contains refresh token that
	// can be used to update access token
	Refresh string `json:"refresh_token,omitempty"`
	// Expiry holds token expiration time
	Expiry time.Time `json:"token_expires,omitempty"`
}

```
What you really need to access user's data with the application that user has authorized is the **access token stored in Access field**.

Token can be saved to and loaded from disk as follows:

```go
err = token.SaveToFile("mytoken.json")
if err != nil {
	fmt.Println("error saving token to file:", err)
}
token, err := OpenAuthToken("mytoken.json")
```

Refresh token stored in Refresh field and can be used to update access token.

```go
// skipping error checks here for brevity
app, _ := yauth.OpenApp("./myapp.json")
token, _ := yauth.OpenAuthToken("./mytoken.json")
newtoken, err := app.Refresh(token.Refresh)
if err == nil {
    newtoken.SaveToFile("mytoken.json")
}
```

