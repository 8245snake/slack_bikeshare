package main

import (
	"fmt"
	"strings"

	bikeshareapi "github.com/8245snake/bikeshare-client"
	"github.com/slack-go/slack"
)

//InteractiveHandler ハンドラ
func InteractiveHandler(message *slack.InteractionCallback) {
	switch message.Type {
	case slack.InteractionTypeShortcut:
		//ショートカット
		fmt.Printf("%v\n", message)
	case slack.InteractionTypeViewSubmission:
		//確定ボタン
		for blockID := range message.View.State.Values {
			//どの画面の確定ボタンかを判別するためにInputブロックのIDを使用
			switch blockID {
			case BlcSearch:
				//駐輪場検索
				SubmitBikeSearch(message)
			case BlcDetail:
				//詳細検索
				SubmitSpotDetail(message)
			}
		}
	case slack.InteractionTypeBlockActions:
		//アクションブロック全般
		for _, blockAction := range message.ActionCallback.BlockActions {
			switch slack.MessageElementType(blockAction.Type) {
			case slack.METButton:
				//ボタン
				switch blockAction.ActionID {
				case BlcSearch:
					//駐輪場検索
					ShowSearchDialog(message)
				case ActDetail:
					//詳細ボタン
					ShowDetailDialog(message, blockAction.Value)
				}
			case slack.METCheckboxGroups:
			case slack.METDatepicker:
			case slack.METImage:
			case slack.METOverflow:
			case slack.METPlainTextInput:
			case slack.METRadioButtons:
			}
		}
	}
}

//ShowSearchDialog 検索ダイアログを出す
func ShowSearchDialog(message *slack.InteractionCallback) {
	view := MakeSearchDialog()
	_, err := api.OpenView(message.TriggerID, view)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

//ShowDetailDialog 検索ダイアログを出す
func ShowDetailDialog(message *slack.InteractionCallback, portCode string) {
	view := MakeDetailDialog(portCode)
	_, err := api.OpenView(message.TriggerID, view)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

//SubmitBikeSearch 検索ボタン押下時
func SubmitBikeSearch(message *slack.InteractionCallback) {
	//検索
	query := message.View.State.Values[BlcSearch][ActSearch].Value
	spots, err := bikeAPI.GetPlaces(bikeshareapi.SearchPlacesOption{Query: query})
	if err != nil {
		errView := MakeErrorView("検索に失敗しました")
		api.PublishView(message.User.ID, errView, "")
		return
	}
	//結果表示
	view := MakeSearchResultView(spots)
	_, err = api.PublishView(message.User.ID, view, "")
	if err != nil {
		text := "結果の表示に失敗しました。"
		if len(spots) > 30 {
			text += fmt.Sprintf("\n検索結果が%d件あります。検索キーワードを絞れば正しく表示される可能性があります。", len(spots))
		}
		errView := MakeErrorView(text)
		api.PublishView(message.User.ID, errView, "")
		return
	}
}

//SubmitSpotDetail 詳細検索結果
func SubmitSpotDetail(message *slack.InteractionCallback) {
	//検索
	data := message.View.State.Values[BlcDetail]
	for _, val := range data {
		var area, spot string
		var days []string
		for _, opt := range val.SelectedOptions {
			if len(days) > 3 {
				break
			}
			valArr := strings.Split(opt.Value, "_")
			if len(valArr) < 2 {
				continue
			}
			area, spot = SplitAreaSpotCode(valArr[0])
			days = append(days, valArr[1])

		}
		view := MakeDetailView(area, spot, days)
		api.PublishView(message.User.ID, view, "")
		return
	}
}
