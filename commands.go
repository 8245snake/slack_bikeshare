package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

//CommandHandler ハンドラ
func CommandHandler(command slack.SlashCommand, w http.ResponseWriter) {

	if command.Command != "/bikeshare" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	switch command.Text {
	case "line login":
		LoginLine(command, w)
	}
}

//LoginLine LINEユーザーと紐付ける
func LoginLine(command slack.SlashCommand, w http.ResponseWriter) {
	rand.Seed(time.Now().UnixNano())
	pin := strconv.Itoa(rand.Int())[:4]
	text := fmt.Sprintf("LINEアカウント <https://lin.ee/8zs3qxk|シェアサイクル台数検索> のトークに以下の４桁のPINコードを送信してください。\n PIN : `%s` ", pin)
	text += "\nPINコードは今から10分間有効です。"
	params := CreateRepryBlockMessage(text)
	b, err := json.Marshal(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	//TODO : DBにセッションを保存する
}

//CreateRepryBlockMessage 返信
func CreateRepryBlockMessage(text string) slack.Message {
	txtBlock := CreateMarkdownText(text)
	messageFrame := slack.NewSectionBlock(nil, []*slack.TextBlockObject{txtBlock}, nil)
	return slack.NewBlockMessage(messageFrame)
}
