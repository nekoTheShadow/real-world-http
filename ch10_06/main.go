package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

//go:embed client_info.json
var accessTokenBytes []byte

type ClientInfo struct {
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"ClientSecret"`
}

func main() {
	var clientInfo ClientInfo
	if err := json.Unmarshal(accessTokenBytes, &clientInfo); err != nil {
		panic(err)
	}

	conf := &oauth2.Config{
		ClientID:     clientInfo.ClientID,
		ClientSecret: clientInfo.ClientSecret,
		Scopes:       []string{"user:email", "gist"},
		Endpoint:     github.Endpoint,
	}
	var token *oauth2.Token

	file, err := os.Open("access_token.json")
	if os.IsNotExist(err) {
		url := conf.AuthCodeURL("your state", oauth2.AccessTypeOnline)
		code := make(chan string)
		var server *http.Server
		server = &http.Server{
			Addr: ":18888",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, "<html><script>widow.open('about:blank', '_self').close()</scipt></html>")
				w.(http.Flusher).Flush()

				code <- r.URL.Query().Get("code")
				server.Shutdown(context.Background())
			}),
		}
		go server.ListenAndServe()
		open.Start(url)
		token, err := conf.Exchange(context.Background(), <-code)
		if err != nil {
			panic(err)
		}
		file, err := os.Create("access_token.json")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(file).Encode(token)
	} else if err == nil {
		token = &oauth2.Token{}
		json.NewDecoder(file).Decode(token)
	} else {
		panic(err)
	}
	client := oauth2.NewClient(context.Background(), conf.TokenSource(context.Background(), token))
	getEmail(client)
	createGist(client)
}

func getEmail(client *http.Client) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	emails, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(emails))
}

func createGist(client *http.Client) {
	gist := `{
		"description": "API example",
		"public": true,
		"files": {
			"hello_from_rest_api.txt": {
				"content": "Hello World!"
			}
		}
	}`
	resp2, err := client.Post("https://api.github.com/gists", "application/json", strings.NewReader(gist))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp2.Status)
	defer resp2.Body.Close()
	type GistResult struct {
		Url string `json:"html_url"`
	}
	gistResult := &GistResult{}
	err = json.NewDecoder(resp2.Body).Decode(&gistResult)
	if err != nil {
		panic(err)
	}
	if gistResult.Url != "" {
		open.Start(gistResult.Url)
	}
}
