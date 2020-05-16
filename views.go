package main

import (
	"fmt"
	"strings"
	"time"

	bikeshareapi "github.com/8245snake/bikeshare-client"
	"github.com/slack-go/slack"
)

//WeekDays 曜日
var WeekDays = [7]string{"日", "月", "火", "水", "木", "金", "土"}

//MakeHomeViewTest 実験
func MakeHomeViewTest() slack.HomeTabViewRequest {
	//テキストセクション
	frame1 := slack.NewSectionBlock(nil, []*slack.TextBlockObject{}, nil)
	frame1.Fields = append(frame1.Fields, slack.NewTextBlockObject("mrkdwn", "*太字*\nなんか文章1", false, false))
	frame1.Fields = append(frame1.Fields, slack.NewTextBlockObject("mrkdwn", "*太字*\nなんか文章2", false, false))
	//画像セクション
	frame2 := slack.NewImageBlock("https://api.slack.com/img/blocks/bkb_template_images/beagle.png", "画像のツールチップ", "",
		slack.NewTextBlockObject("plain_text", "画像のタイトル", false, false))
	//コンテキスト
	frame3 := slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "コンテキスト文章", false, false))
	//仕切り
	divider := slack.NewDividerBlock()
	//Selectボックスいろいろ（ActionBlockの中に複数コントロールを入れる）
	frame4 := slack.NewActionBlock("")
	frame4.Elements.ElementSet = append(frame4.Elements.ElementSet,
		slack.NewOptionsSelectBlockElement("conversations_select", slack.NewTextBlockObject("plain_text", "会話", false, false), ""))
	frame4.Elements.ElementSet = append(frame4.Elements.ElementSet,
		slack.NewOptionsSelectBlockElement("channels_select", slack.NewTextBlockObject("plain_text", "チャンネル", false, false), ""))
	frame4.Elements.ElementSet = append(frame4.Elements.ElementSet,
		slack.NewOptionsSelectBlockElement("users_select", slack.NewTextBlockObject("plain_text", "ユーザー", false, false), ""))
	elem := slack.NewOptionsSelectBlockElement("static_select", slack.NewTextBlockObject("plain_text", "自由定義項目", false, false), "")
	elem.Options = append(elem.Options, slack.NewOptionBlockObject("item_1", slack.NewTextBlockObject("plain_text", "Item1", false, false)))
	elem.Options = append(elem.Options, slack.NewOptionBlockObject("item_2", slack.NewTextBlockObject("plain_text", "Item2", false, false)))
	elem.Options = append(elem.Options, slack.NewOptionBlockObject("item_3", slack.NewTextBlockObject("plain_text", "Item3", false, false)))
	frame4.Elements.ElementSet = append(frame4.Elements.ElementSet, elem)
	//セクション内に画像
	frame5 := slack.NewSectionBlock(nil, []*slack.TextBlockObject{}, nil)
	frame5.Fields = append(frame5.Fields, slack.NewTextBlockObject("mrkdwn", "*太字*\n画像の説明など", false, false))
	frame5.Accessory = slack.NewAccessory(slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/palmtree.png", "ツールチップです"))
	//セクション内にボタン
	frame6 := slack.NewSectionBlock(nil, []*slack.TextBlockObject{}, nil)
	frame6.Fields = append(frame6.Fields, slack.NewTextBlockObject("mrkdwn", "*ボタンの説明など*", false, false))
	frame6.Accessory = slack.NewAccessory(slack.NewButtonBlockElement("", "button_pushed", slack.NewTextBlockObject("plain_text", "押してみよう", false, false)))
	//セクション内にSelectボックス
	frame7 := slack.NewSectionBlock(nil, []*slack.TextBlockObject{}, nil)
	frame7.Fields = append(frame7.Fields, slack.NewTextBlockObject("mrkdwn", "*Selectボックスの説明など*", false, false))
	frame7.Accessory = slack.NewAccessory(elem)
	//セクション内にマルチSelectボックス
	frame8 := slack.NewSectionBlock(nil, []*slack.TextBlockObject{}, nil)
	frame8.Fields = append(frame8.Fields, slack.NewTextBlockObject("mrkdwn", "*マルチSelectボックスの説明など*", false, false))
	multi := slack.NewOptionsMultiSelectBlockElement("multi_static_select", slack.NewTextBlockObject("plain_text", "自由定義項目（複数）", false, false), "")
	multi.Options = append(multi.Options, slack.NewOptionBlockObject("item_1", slack.NewTextBlockObject("plain_text", "Item1", false, false)))
	multi.Options = append(multi.Options, slack.NewOptionBlockObject("item_2", slack.NewTextBlockObject("plain_text", "Item2", false, false)))
	multi.Options = append(multi.Options, slack.NewOptionBlockObject("item_3", slack.NewTextBlockObject("plain_text", "Item3", false, false)))
	frame8.Accessory = slack.NewAccessory(multi)

	blockSet := []slack.Block{frame1, frame2, frame3, divider, frame4, frame5, frame6, frame7, frame8}
	view := slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: slack.Blocks{BlockSet: blockSet},
	}
	return view
}

//CreateSearchFrame 検索ボタンなどが乗ったブロック作成
func CreateSearchFrame() *slack.ActionBlock {
	button := slack.NewButtonBlockElement(BlcSearch, BtnOpenSearchDialog, CreatePlainText("駐輪場検索"))
	frame := slack.NewActionBlock("")
	frame.Elements.ElementSet = append(frame.Elements.ElementSet, button)
	return frame
}

//MakeHomeView ホーム画面
func MakeHomeView() slack.HomeTabViewRequest {
	//TODO : ユーザーごとに画面を変えたい
	block := CreateSearchFrame()
	blockSet := []slack.Block{block}
	view := slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: slack.Blocks{BlockSet: blockSet},
	}
	return view
}

//MakeSearchDialog 検索画面作成
func MakeSearchDialog() slack.ModalViewRequest {
	//検索欄
	inputelement := slack.NewPlainTextInputBlockElement(CreatePlainText("検索クエリ（例：'A1', 'A2-03', '都庁'）"), ActSearch)
	inputblock := slack.NewInputBlock(BlcSearch, CreatePlainText("自由検索"), inputelement)
	//リクエスト
	view := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  CreatePlainText("駐輪場検索"),
		Blocks: slack.Blocks{BlockSet: []slack.Block{inputblock, slack.NewDividerBlock()}},
		Submit: CreatePlainText("検索"),
		Close:  CreatePlainText("閉じる"),
	}
	return view
}

//MakeDetailDialog 検索画面作成
func MakeDetailDialog(portCode string) slack.ModalViewRequest {
	multi := slack.NewOptionsMultiSelectBlockElement("multi_static_select", slack.NewTextBlockObject("plain_text", "日付選択（複数）", false, false), "")
	day := time.Now()
	for i := 0; i < 30; i++ {
		value := portCode + "_" + day.Format("20060102")
		multi.Options = append(multi.Options, slack.NewOptionBlockObject(value, CreatePlainText(fmt.Sprintf("%s (%s)", day.Format("2006/01/02"), WeekDays[day.Weekday()]))))
		day = day.AddDate(0, 0, -1)
	}
	inputBlock := slack.NewInputBlock(BlcDetail, CreatePlainText("グラフを表示する日付を選択してください（最大３つ）"), multi)
	//画像セクション
	area, spot := SplitAreaSpotCode(portCode)
	graph, _ := bikeAPI.GetGraph(
		bikeshareapi.SearchGraphOption{
			Area:      area,
			Spot:      spot,
			DrawTitle: true,
			Days:      []string{time.Now().Format("20060102")},
		})
	pictureBlock := slack.NewImageBlock(graph.URL, "画像のツールチップ", "",
		CreatePlainText("台数の経時変化"))
	//リクエスト
	view := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  CreatePlainText("駐輪場検索"),
		Blocks: slack.Blocks{BlockSet: []slack.Block{inputBlock, pictureBlock}},
		Submit: CreatePlainText("検索"),
		Close:  CreatePlainText("閉じる"),
	}
	return view
}

//MakeDetailView 検索結果画面
func MakeDetailView(area string, spot string, days []string) slack.HomeTabViewRequest {
	searchBlock := CreateSearchFrame()
	//画像セクション
	graph, _ := bikeAPI.GetGraph(
		bikeshareapi.SearchGraphOption{
			Area:      area,
			Spot:      spot,
			DrawTitle: true,
			Days:      days,
		})
	pictureBlock := slack.NewImageBlock(graph.URL, "画像のツールチップ", "",
		CreatePlainText("台数の経時変化"))
	blockSet := []slack.Block{searchBlock, pictureBlock}
	view := slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: slack.Blocks{BlockSet: blockSet},
	}
	return view
}

//MakeSearchResultView 検索結果画面
func MakeSearchResultView(spotinfoList []bikeshareapi.SpotInfo) slack.HomeTabViewRequest {
	block := CreateSearchFrame()
	blockSet := []slack.Block{block, slack.NewDividerBlock()}
	for _, info := range spotinfoList {
		text := fmt.Sprintf("*[%s-%s] %s*\n", info.Area, info.Spot, info.Name)
		context := ""
		if len(info.Counts) > 0 {
			text += fmt.Sprintf("%d台", info.Counts[0].Count)
			context = fmt.Sprintf("%s", info.Counts[0].Time.Format("2006/01/02 15:04"))
		} else {
			text += "台数が取得できませんでした"
		}

		txtBlock := CreateMarkdownText(text)
		frame := slack.NewSectionBlock(nil, []*slack.TextBlockObject{txtBlock}, nil)
		//詳細ボタン
		button := slack.NewAccessory(slack.NewButtonBlockElement(ActDetail, info.Area+"-"+info.Spot, CreatePlainText("詳細")))
		frame.Accessory = button
		blockSet = append(blockSet, frame)
		if context != "" {
			blockSet = append(blockSet, slack.NewContextBlock("", CreateMarkdownText(context)))
		}
		blockSet = append(blockSet, slack.NewDividerBlock())
	}

	view := slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: slack.Blocks{BlockSet: blockSet},
	}
	return view
}

//MakeErrorView エラー画面作成
func MakeErrorView(message string) slack.HomeTabViewRequest {
	searchFrame := CreateSearchFrame()
	txtBlock := CreateMarkdownText(message)
	messageFrame := slack.NewSectionBlock(nil, []*slack.TextBlockObject{txtBlock}, nil)
	blockSet := []slack.Block{searchFrame, messageFrame}
	view := slack.HomeTabViewRequest{
		Type:   slack.VTHomeTab,
		Blocks: slack.Blocks{BlockSet: blockSet},
	}
	return view
}

//CreatePlainText プレーンテキスト
func CreatePlainText(text string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("plain_text", text, false, false)
}

//CreateMarkdownText マークダウン
func CreateMarkdownText(text string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("mrkdwn", text, false, false)
}

//SplitAreaSpotCode コードをareaとspotに分離
func SplitAreaSpotCode(portCode string) (area, spot string) {
	arr := strings.Split(portCode, "-")
	if len(arr) < 2 {
		return
	}
	area = arr[0]
	spot = arr[1]
	return
}
