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
	"strings"
	"time"

	"github.com/demisto/slack"

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
		Scopes:      []string{"chat:write:bot", "incoming-webhook", "commands", "bot"},
		RedirectURL: "https://festivus.nivas.hr/auth",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/authorize",
			TokenURL: "https://slack.com/api/oauth.access", // not actually used here
		},
	}
	url := conf.AuthCodeURL(globalState.auth)
	http.Redirect(w, r, url, http.StatusFound)
}

var (
	s             *slack.Slack
	info          *slack.RTMStartReply // The global info for the team
	currChannelID string               // The ID of the channel
	files         []slack.File         // The files for the team
)

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
	token, err := slack.OAuthAccess(*clientID, *clientSecret, code, "")
	if err != nil {
		writeError(w, 401, err.Error())
		return
	}
	s, err = slack.New(slack.SetToken(token.AccessToken))
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	// Get our own user id
	test, err := s.AuthTest()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.Write([]byte(fmt.Sprintf("OAuth successful for team %s and user %s", test.Team, test.User)))
	log.Printf("%#v", test)

	in := make(chan *slack.Message)
	info, err = s.RTMStart("Your URL", in, nil)

	currChannelID = channelID("gophergala")
	postMessage("Hello world!")

}

func postMessage(msg string) {
	m := &slack.PostMessageRequest{
		AsUser:  false, // bot is not user? user should have client special scope ? xxx todo
		Channel: currChannelID,
		Text:    msg,
	}
	_, err := s.PostMessage(m, true)
	if err != nil {
		fmt.Printf("Unable to post to channel %s - %v\n", channelName(currChannelID), err)
	}
}
func channelID(ch string) string {
	// First, let's see if the given ch is actually already an ID
	name := channelName(ch)
	if name != "" {
		return ch
	}
	for i := range info.Channels {
		if strings.ToLower(info.Channels[i].Name) == strings.ToLower(ch) {
			return info.Channels[i].ID
		}
	}
	for i := range info.Groups {
		if strings.ToLower(info.Groups[i].Name) == strings.ToLower(ch) {
			return info.Groups[i].ID
		}
	}
	for i := range info.IMS {
		if strings.ToLower(userNameByID(info.IMS[i].User)) == strings.ToLower(ch) {
			return info.IMS[i].ID
		}
	}
	return ""
}
func channelName(ch string) string {
	if ch == "" {
		return ""
	}
	switch ch[0] {
	case 'C':
		for i := range info.Channels {
			if info.Channels[i].ID == ch {
				return info.Channels[i].Name
			}
		}
	case 'G':
		for i := range info.Groups {
			if info.Groups[i].ID == ch {
				return info.Groups[i].Name
			}
		}
	case 'D':
		for i := range info.IMS {
			if info.IMS[i].ID == ch {
				return userNameByID(info.IMS[i].User)
			}
		}
	}
	return ""
}

// userNameByID if user is not found then just use ID
func userNameByID(id string) string {
	uname := id
	u := findUser(id)
	if u != nil {
		uname = u.Name
	}
	return uname
}
func findUser(id string) *slack.User {
	for i := range info.Users {
		if info.Users[i].ID == id {
			return &info.Users[i]
		}
	}
	return nil
}

func switchChannel(ch string) bool {
	id := channelID(ch)
	if id != "" {
		currChannelID = id
		return true
	}
	return false
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
