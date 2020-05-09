package main

import (
	"fmt"

	bikeshareapi "github.com/8245snake/bikeshare-client"
	"github.com/slack-go/slack"
)

//ShowSearchDialog 検索ダイアログを出す
func ShowSearchDialog(message slack.InteractionCallback) {
	view := MakeSearchDialog()
	_, err := api.OpenView(message.TriggerID, view)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

//SubmitBikeSearch 検索ボタン押下時
func SubmitBikeSearch(message slack.InteractionCallback) {
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
