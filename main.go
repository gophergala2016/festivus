package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	// "github.com/demisto/slack"
	"github.com/gophergala2016/festivus/holidays"
	"github.com/nlopes/slack"

	"golang.org/x/oauth2"
)

var (
	address      = flag.String("address", ":2016", "Which address should I listen on")
	clientID     = flag.String("client_id", "", "The client ID from https://api.slack.com/applications")
	clientSecret = flag.String("client_secret", "", "The client secret from https://api.slack.com/applications")

	CharsetUTF8                = "charset=utf-8"
	ApplicationJSON            = "application/json"
	ApplicationJSONCharsetUTF8 = ApplicationJSON + "; " + CharsetUTF8
	ContentType                = "Content-Type"
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
		// scope=incoming-webhook,commands,bot
		// incoming-webhook - post from your app to a single Slack channel.
		// za rtm kazu da treba rtm:stream pa kazu da je unknown..hmm, probam sa client.
		// Scopes:      []string{"commands", "bot", "chat:write:bot", "client"}, // za ovo javlja da mixam depresiated scopove argh...
		Scopes:      []string{"incoming-webhook", "commands"}, // jel moguce da je "samo" "client" dovoljno? | maknuto "bot", "chat:write:bot"
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
		// writeError(w, 401, err.Error())
		writeError(w, 500, err.Error())
		return
	}
	log.Println("got:", token, scope)

	api := slack.New(token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)

	// params := slack.PostMessageParameters{}
	// attachment := slack.Attachment{
	// 	AuthorName: "FestivusBot",
	// 	Pretext:    "some pretext",
	// 	Text:       "some text",
	// 	// Uncomment the following part to send a field too
	// 	/*
	// 		Fields: []slack.AttachmentField{
	// 			slack.AttachmentField{
	// 				Title: "a",
	// 				Value: "no",
	// 			},
	// 		},
	// 	*/
	// }
	// params.Attachments = []slack.Attachment{attachment}
	// channelID, timestamp, err := api.PostMessage("C0JAT1RHS", "Hello world", params)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }
	// fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	w.Write([]byte(fmt.Sprintf("Install done.")))

	return

	// rtm := api.NewRTM()
	// 	go rtm.ManageConnection()

	// Loop:
	// 	for {
	// 		select {
	// 		case msg := <-rtm.IncomingEvents:
	// 			fmt.Print("Event Received: ")
	// 			switch ev := msg.Data.(type) {
	// 			case *slack.HelloEvent:
	// 				// Ignore hello

	// 			case *slack.ConnectedEvent:
	// 				fmt.Println("Infos:", ev.Info)
	// 				fmt.Println("Connection counter:", ev.ConnectionCount)
	// 				// Replace #general with your Channel ID
	// 				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

	// 			case *slack.MessageEvent:
	// 				fmt.Printf("Message: %v\n", ev)

	// 			case *slack.PresenceChangeEvent:
	// 				fmt.Printf("Presence Change: %v\n", ev)

	// 			case *slack.LatencyReport:
	// 				fmt.Printf("Current latency: %v\n", ev.Value)

	// 			case *slack.RTMError:
	// 				fmt.Printf("Error: %s\n", ev.Error())

	// 			case *slack.InvalidAuthEvent:
	// 				fmt.Printf("Invalid credentials")
	// 				break Loop

	// 			default:

	// 				// Ignore other events..
	// 				// fmt.Printf("Unexpected: %v\n", msg.Data)
	// 			}
	// 		}
	// 	}

}

// RADI od prve, ali baca missing_scope","needed":"groups:read",
// groups, err := api.GetGroups(false)
// if err != nil {
// 	fmt.Printf("%s\n", err)
// 	return
// }
// for _, group := range groups {
// 	fmt.Printf("Id: %s, Name: %s\n", group.ID, group.Name)
// }

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

// home displays the add-to-slack button
func home(w http.ResponseWriter, r *http.Request) {
	slackbutton := `<img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x">`
	w.Write([]byte(`<html><head><title>Slack OAuth Test</title></head><body><a href="/add">` + slackbutton + `</a></body></html>`))
}

// https://festivus.nivas.hr/slack/festivus
func festivusCmd(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte(`hello world`))

	// response_type is in_channel, both the response message and the initial message typed by the user will be shared in the channel
	// response_type to ephemeral is the same as not including the response type at all, and the response message will be visible only to the user that issued the command
	// {
	//     "response_type": "in_channel",
	//     "text": "It's 80 degrees right now.",
	//     "attachments": [
	//         {
	//             "text":"Partly cloudy today and tomorrow"
	//         }
	//     ]
	// }

	params := r.FormValue("text")

	// empty call means "days till festivus"
	if len(params) == 0 {
		daysTillFestivus := Festivus(time.Now())
		festivusDay := FestivusDate(time.Now())

		err := JSON(
			w,
			http.StatusOK,
			struct {
				Text string `json:"text"`
			}{
				fmt.Sprintf("%d days til Festivus (%s).", daysTillFestivus, festivusDay.Format("02.01.2006.")),
			})

		if err != nil {
			writeError(w, 500, err.Error())
		}

		return
	}

	countryCode := strings.Trim(params, " ")

	curDir, err := os.Getwd()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	calPath := curDir + "/calendars/" + countryCode + ".txt"
	log.Printf("cal path %q", calPath)

	if _, err := os.Stat(calPath); err != nil && os.IsNotExist(err) {
		err := JSON(
			w,
			http.StatusOK,
			struct {
				Text string `json:"text"`
			}{
				fmt.Sprintf("Unsupported country %q", countryCode),
			})

		if err != nil {
			writeError(w, 500, err.Error())
		}

		return
	}

	days, err := holidays.New(countryCode, curDir+"/calendars/")
	if err != nil {
		if err != nil {
			writeError(w, 500, err.Error())
		}

		return
	}

	var sDays string
	for _, d := range days {
		sDays = fmt.Sprintf("%s\n%s", sDays, d.String())
	}

	err = JSON(
		w,
		http.StatusOK,
		struct {
			Text string `json:"text"`
		}{
			fmt.Sprintf("%s", sDays),
		})

	if err != nil {
		writeError(w, 500, err.Error())
	}

	return
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
	http.HandleFunc("/slack/festivus", festivusCmd)
	http.HandleFunc("/add", addToSlack)
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*address, nil))
}

func JSON(w http.ResponseWriter, code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	w.Write(b)
	return
}

// DaysBetween returns days between dates.
func DaysBetween(from, to time.Time) int {
	// convert diff hours to days
	d := to.Sub(from).Hours() / 24
	return int(math.Abs(d))
}

// Festivus returns number of days from to today to festivus
func Festivus(today time.Time) int {
	festDate := FestivusDate(today)
	return DaysBetween(today, festDate)
}

// FestivusDate
func FestivusDate(today time.Time) (festivus time.Time) {
	year := time.Now().Year()

	festivus = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	if today.After(festivus) {
		year++
		festivus = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	}

	return
}
