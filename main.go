package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	webServiceBaseURL = "https://example.com/api"
)

func main() {
	webServiceBaseURL = os.Getenv("WEBSERVICE_URL")
	botAPIKey := os.Getenv("TELEGRAM_APIKEY")
	groupID, err := strconv.ParseInt(os.Getenv("TELEGRAM_GROUPID"), 10, 64)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(botAPIKey)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot 名称: %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	go factorioMessageBridge(bot, groupID)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("<%s> %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			handleCommands(bot, update)
			continue
		}

		if update.Message.Chat.ID == groupID {
			sendMessage(update.Message.From.UserName, update.Message.Text)
		}
	}
}

func handleCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyToMessageID = update.Message.MessageID

	switch update.Message.Command() {
	case "help":
		msg.Text = "type /factorioplayers"
	case "factorioplayers":
		players, err := getOnlinePlayers()
		if err != nil {
			msg.Text = err.Error()
		} else {
			if len(*players) == 0 {
				msg.Text = "好像并没有玩家在游戏中……"
			} else {
				for _, p := range *players {
					msg.Text += p.Name + ", "
				}
				msg.Text = strings.TrimSuffix(msg.Text, ", ")
			}
		}
	}

	if msg.Text == "" {
		return
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func factorioMessageBridge(bot *tgbotapi.BotAPI, groupID int64) {
	for {
		u, err := getUpdate()
		if err != nil {
			log.Println("获取 Factorio 事件更新失败: ", err)
		}
		switch u.Type {
		case "empty":
		case "console-me":
			message := fmt.Sprintf("[Factorio] *%s %s\n", u.PlayerName, u.Message)
			log.Println(message)
			msg := tgbotapi.NewMessage(groupID, message)
			bot.Send(msg)
			continue
		case "console-chat":
			message := fmt.Sprintf("[Factorio] <%s> %s\n", u.PlayerName, u.Message)
			log.Println(message)
			msg := tgbotapi.NewMessage(groupID, message)
			bot.Send(msg)
			continue
		case "oops":
			message := fmt.Sprintf("[Factorio] oops, %s\n", u.PlayerName)
			log.Println(message)
			msg := tgbotapi.NewMessage(groupID, message)
			bot.Send(msg)
			message = fmt.Sprintf("[Factorio] %s\n", u.Sarcasm)
			log.Println(message)
			msg = tgbotapi.NewMessage(groupID, message)
			bot.Send(msg)
			continue
		default:
			log.Println("Unknown update type:", u.Type)
		}

		time.Sleep(time.Second * 5)
	}
}
