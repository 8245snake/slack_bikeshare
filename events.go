package main

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

//EventHandler ハンドラ
func EventHandler(eventsAPIEvent *slackevents.EventsAPIEvent) {
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, ok := innerEvent.Data.(*slackevents.AppMentionEvent)
			if !ok {
				return
			}
			api.PostMessage(ev.Channel, slack.MsgOptionText("AppMentionEvent", false))
		case *slackevents.MessageEvent:
			_, ok := innerEvent.Data.(*slackevents.MessageEvent)
			if !ok {
				return
			}
			api.PostMessage(ev.Channel, slack.MsgOptionText("MessageEvent", false))
		case *slackevents.EventsAPICallbackEvent:
			_, ok := innerEvent.Data.(*slackevents.MessageEvent)
			if !ok {
				return
			}
			fmt.Println("EventsAPICallbackEvent")
		case *slackevents.AppHomeOpenedEvent:
			if s, ok := innerEvent.Data.(*slackevents.AppHomeOpenedEvent); ok {
				AppHomeOpened(s)
			}
		}

	}

}

//AppHomeOpened ホーム画面が開かれたとき
func AppHomeOpened(event *slackevents.AppHomeOpenedEvent) {
	view := MakeHomeView()
	api.PublishView(event.User, view, "")
}
