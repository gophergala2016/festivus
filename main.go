package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	// "github.com/demisto/slack"
	"github.com/nlopes/slack"

	"golang.org/x/oauth2"
)

var (
	address      = flag.String("address", ":2016", "Which address should I listen on")
	clientID     = flag.String("client_id", "", "The client ID from https://api.slack.com/applications")
	clientSecret = flag.String("client_secret", "", "The client secret from https://api.slack.com/applications")
)

type state struct {
	auth string
	ts   time.Time
}

// globalState is an example of how to store a state between calls
var globalState state

// writeError writes an error to the reply - example only
func writeError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(err))
	log.Println("writeError", err)
}

// addToSlack initializes the oauth process and redirects to Slack
func addToSlack(w http.ResponseWriter, r *http.Request) {
	// Just generate random state
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		writeError(w, 500, err.Error())
	}
	globalState = state{auth: hex.EncodeToString(b), ts: time.Now()}
	conf := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		// Scopes:       []string{"client"}, // special scope (!) - Allows applications to connect to slack as a client, and post messages on behalf of the user.
		// incoming-webhook - post from your app to a single Slack channel.
		Scopes:      []string{"commands", "bot"},
		RedirectURL: "https://festivus.nivas.hr/auth",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access", // not actually used here
		},
	}
	url := conf.AuthCodeURL(globalState.auth)
	http.Redirect(w, r, url, http.StatusFound)
}

// auth receives the callback from Slack, validates and displays the user information
func auth(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	errStr := r.FormValue("error")
	if errStr != "" {
		writeError(w, 401, errStr)
		return
	}
	if state == "" || code == "" {
		writeError(w, 400, "Missing state or code")
		return
	}
	if state != globalState.auth {
		writeError(w, 403, "State does not match")
		return
	}
	// As an example, we allow only 5 min between requests
	if time.Since(globalState.ts) > 5*time.Minute {
		writeError(w, 403, "State is too old")
		return
	}

	token, scope, err := slack.GetOAuthToken(*clientID, *clientSecret, code, "", true)
	if err != nil {
		writeError(w, 401, err.Error())
		return
	}
	log.Println("got:", token, scope)

	api := slack.New(token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)
	groups, err := api.GetGroups(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("Id: %s, Name: %s\n", group.ID, group.Name)
	}

	// token, err := slack.OAuthAccess(*clientID, *clientSecret, code, "")
	// if err != nil {
	// 	writeError(w, 401, err.Error())
	// 	return
	// }
	// s, err := slack.New(slack.SetToken(token.AccessToken))
	// if err != nil {
	// 	writeError(w, 500, err.Error())
	// 	return
	// }
	// // Get our own user id
	// test, err := s.AuthTest()
	// if err != nil {
	// 	writeError(w, 500, err.Error())
	// 	return
	// }
	// w.Write([]byte(fmt.Sprintf("OAuth successful for team %s and user %s", test.Team, test.User)))
	// log.Printf("%#v", test)
}

// home displays the add-to-slack button
func home(w http.ResponseWriter, r *http.Request) {
	slackbutton := `<img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x">`
	w.Write([]byte(`<html><head><title>Slack OAuth Test</title></head><body><a href="/add">` + slackbutton + `</a></body></html>`))
}

func main() {
	flag.Parse()
	if *clientID == "" || *clientSecret == "" {
		fmt.Print(`
You must specify the client ID and client secret from https://api.slack.com/applications
Usage: ./festivus --address ":2016" --client_id "YOUR_ID" --client_secret "YOUR_SECRET"

`)
		os.Exit(1)
	}
	http.HandleFunc("/add", addToSlack)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*address, nil))
}

// DaysBetween returns days between dates.
func DaysBetween(from, to time.Time) int {
	// convert diff hours to days
	d := to.Sub(from).Hours() / 24
	return int(math.Abs(d))
}

// Festivus returns number of days from to today to festivus
func Festivus(today time.Time) int {

	year := time.Now().Year()

	festDate := time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	if today.After(festDate) {
		year++
		festDate = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	}
	return DaysBetween(today, festDate)
}
