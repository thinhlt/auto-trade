package bot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"anidiot.com/auto-trade/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Bot struct {
	sigChan  chan string
	timer    <-chan time.Time
	stopChan chan string
	wg       *sync.WaitGroup
}

var channel *tgbotapi.BotAPI
var messageChan chan string
var donC chan struct{}
var mutex = &sync.Mutex{}
var warningMap = make(map[string]struct{}, 0)

func (bot Bot) initBotChannel(token string, channelID int64) {
	fmt.Println("init bot with config", token, channelID)
	telegram, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
	}
	// telegram.Debug = true
	bot.wg.Add(1)

	log.Printf("Authorized on account %s", telegram.Self.UserName)
	go func() {
		defer bot.wg.Done()
		for {
			select {
			case signal := <-bot.stopChan:
				if signal == "channel" {
					fmt.Println("stop channel")
					return
				} else {
					bot.stopChan <- signal
				}
			case s := <-bot.sigChan:
				fmt.Println("send message", s)
				message := tgbotapi.NewMessage(channelID, s)
				_, err := telegram.Send(message)
				if err != nil {
					fmt.Println("error when send telegram", err)
				}
				break
			}
		}
	}()
}

type ChatBot struct {
}

func InitBot(warningMap map[string]struct{}) ChatBot {
	token := viper.GetString("rhino.token")
	channelID := viper.GetInt64("telegram.id")
	warningChannelID := viper.GetInt64("telegram.warning_id")
	log.Logger.Info("init bot with config", zap.String("token", token), zap.Int64("channelID", channelID))
	var err error
	channel, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Logger.Error("init bot", zap.Error(err))
	}
	channel.Debug = false
	messageChan = make(chan string, 100)
	donC = make(chan struct{})
	go func(warnMap map[string]struct{}) {
		for {
			select {
			case <-donC:
				log.Logger.Info("Stop signal channel")
				return
			case s := <-messageChan:
				msgPart := strings.Split(s, " ")
				name := msgPart[0]
				if _, present := warnMap[name]; present {
					message := tgbotapi.NewMessage(warningChannelID, s)
					_, err := channel.Send(message)
					if err != nil {
						log.Logger.Error("error when send telegram", zap.Error(err))
						errorPart := strings.Split(err.Error(), "retry after ")
						if len(errorPart) >= 1 {
							time.Sleep(time.Second * 15)
						}
					}
				}
				continue
			}
		}
	}(warningMap)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := channel.GetUpdatesChan(u)

	go func(warnMap map[string]struct{}) {
		for {
			select {
			case <-donC:
				log.Logger.Info("Stop warning channel")
				return
			case update := <-updates:
				if update.Message == nil {
					continue
				}

				log.Logger.Info("command", zap.String("username", update.Message.From.UserName), zap.String("mess", update.Message.Text))
				msgPart := strings.Split(update.Message.Text, ":")
				if len(msgPart) < 2 || (!strings.EqualFold(msgPart[0], "start") && !strings.EqualFold(msgPart[0], "end")) {
					return
				}
				if strings.EqualFold(msgPart[0], "start") {
					warnMap[strings.ToUpper(msgPart[1])] = struct{}{}
				} else {
					if _, ok := warningMap[strings.ToUpper(msgPart[1])]; ok {
						delete(warnMap, strings.ToUpper(msgPart[1]))
					}
				}
			}
		}
	}(warningMap)

	return ChatBot{}
}

func (ChatBot) Close() {
	donC <- struct{}{}
	donC <- struct{}{}
	donC <- struct{}{}
}

func (ChatBot) SendMessage(messages []string, channelID int64) {
	if len(messages) == 0 {
		return
	}
	for _, m := range messages {
		messageChan <- m
	}
}
