package main

import (
	"flag"
	"fmt"

	"github.com/dmfed/yauth"
)

func main() {
	var (
		id      = flag.String("id", "", "your Yandex OAuth app id")
		secret  = flag.String("secret", "", "your Yandex OAuth app secret")
		refresh = flag.String("refresh", "", "refresh token to use (if empty will request user authorization)")
		tokfile = flag.String("f", "", "filename to save token to")
		appfile = flag.String("app", "", "filename with your app credentials")
	)
	flag.Parse()
	var app yauth.Application
	var err error
	if *appfile != "" {
		app, err = yauth.OpenApp(*appfile)
	} else {
		app, err = yauth.NewApp(*id, *secret)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	var token yauth.AuthToken
	if *refresh != "" {
		token, err = app.Refresh(*refresh)
	} else {
		token, err = app.RequestUserAuthorization()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	if *tokfile != "" {
		err = token.SaveToFile(*tokfile)
		if err != nil {
			fmt.Println("error saving token to file:", err)
		}
	}
	fmt.Printf("New token:\n%v\n", token)
	if err != nil {
		fmt.Println(err)
		return
	}
}
