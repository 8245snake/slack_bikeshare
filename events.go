package main

import (
	"github.com/slack-go/slack/slackevents"
)

//AppHomeOpened ホーム画面が開かれたとき
func AppHomeOpened(event *slackevents.AppHomeOpenedEvent) {
	view := MakeHomeView()
	api.PublishView(event.User, view, "")
}
