# Yandex authentication with OAuth API 
This library implements one of ways to obtain access token used to access user data in various Yandex services (for example to access Disk REST API) from Yandex OAuth API.

A detailed description of implemented way to obtain access token is available here: https://yandex.com/dev/oauth/doc/dg/reference/simple-input-client.html#simple-input-client

Interaction with user is handled in terminal, so the library is suitable for use in console-only applications. 

First you need to register your application at https://oauth.yandex.ru to receive credentials and define
default scope of permissions that will be requested from user.

## Usage
```go
import "github.com/dmfed/yauth"
```

Use the credentials received when registering your app to create yauth.Application as follows:

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
The app then prints instructions for user to Stdout and starts polling 
Yandex OAuth API to check whether the user has authorized your application. 

To authorize you app the user needs to visit https://ya.ru/device and enter the device code that was printed in terminal. 
The actual link printed to user may differ, because Yandex API returns the URL that needs to be visited. 

If authorization was successful the app returns AuthToken.

## Token

The Token is defined as follows:

```go
// Token holds all values nececcary to access
// Yandex services using Yandex OAuth.
type Token struct {
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
token, err := yauth.OpenToken("mytoken.json")
```

Refresh token value stored in Refresh field and can be used to update access token.

```go
// skipping error checks here for brevity
app, _ := yauth.OpenApp("./myapp.json")
token, _ := yauth.OpenToken("./mytoken.json")
newtoken, err := app.Refresh(token.Refresh)
if err == nil {
    newtoken.SaveToFile("mytoken.json")
}
```

Example cli appcan be found in **example** directory of the repository.

Any suggestions are most welcome. 

