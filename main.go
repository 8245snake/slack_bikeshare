package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	bikeshareapi "github.com/8245snake/bikeshare-client"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

//環境変数
var (
	BotToken    = os.Getenv("SLACK_TOKEN")
	VerifyToken = os.Getenv("SLACK_VERIFY")
	Secretoken  = os.Getenv("SLACK_SECRET")
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New(BotToken)

//bikeAPI APIクライアント
var bikeAPI = bikeshareapi.NewApiClient()

//EventEndpoint イベントPOSTハンドラ
func EventEndpoint(w http.ResponseWriter, r *http.Request) {
	retry := r.Header.Get("X-Slack-Retry-Num")
	if retry != "" {
		log.Printf("X-Slack-Retry-Num = %s", retry)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: VerifyToken}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	go EventHandler(&eventsAPIEvent)
	w.WriteHeader(http.StatusOK)
}

//InteractiveEndpoint イベントPOSTハンドラ
func InteractiveEndpoint(w http.ResponseWriter, r *http.Request) {
	retry := r.Header.Get("X-Slack-Retry-Num")
	if retry != "" {
		log.Printf("X-Slack-Retry-Num = %s", retry)
		return
	}
	// Read request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] Fail to read request body: %v", err)
		return
	}

	// Parse request body
	str, _ := url.QueryUnescape(string(body))
	str = strings.Replace(str, "payload=", "", 1)
	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(str), &message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("[ERROR] Fail to unmarchal json: %v", err)
		return
	}
	// fmt.Println(str)
	go InteractiveHandler(&message)
	w.WriteHeader(http.StatusOK)
}

//CommandEndpoint イベントPOSTハンドラ
func CommandEndpoint(w http.ResponseWriter, r *http.Request) {
	retry := r.Header.Get("X-Slack-Retry-Num")
	if retry != "" {
		log.Printf("X-Slack-Retry-Num = %s", retry)
		return
	}
	verifier, err := slack.NewSecretsVerifier(r.Header, Secretoken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	CommandHandler(command, w)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/events-endpoint", EventEndpoint)
	http.HandleFunc("/interactive-endpoint", InteractiveEndpoint)
	http.HandleFunc("/command", CommandEndpoint)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":5050", nil)
}
