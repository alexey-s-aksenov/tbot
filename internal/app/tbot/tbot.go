package tbot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexey-s-aksenov/tbot/internal/pkg/deluge"
	"github.com/alexey-s-aksenov/tbot/internal/pkg/joke"
	"github.com/alexey-s-aksenov/tbot/internal/pkg/saveusers"
	"github.com/alexey-s-aksenov/tbot/internal/pkg/schedule"
	"github.com/alexey-s-aksenov/tbot/internal/pkg/weather"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var myKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Шутка"),
		tgbotapi.NewKeyboardButton("Торрент"),
	),
)

var regUserKeyb = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Шутка"),
		tgbotapi.NewKeyboardButton("Отменить подписку"),
	),
)

var newUserKeyb = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Шутка"),
		tgbotapi.NewKeyboardButton("Подписка"),
	),
)

// CreateBot function loads configs and starts bot
func CreateBot() {

	// Intialization
	// 1. get logger instance
	logger := newFileLogger()

	// 2. load config
	config, err := readConfig()
	if err != nil {
		logger.Panic("failed to load bot configuration", err)
	}

	// 3. load registered users and setup goroutin to track new users
	var userConf saveusers.EncConfig
	userConf.File = config.Users.File
	users, err := userConf.Load()
	regUsers := *users

	if err != nil {
		logger.Panic("failed to load registered users", err)
	}

	userChange := make(chan *map[string]int64)
	go func() {
		for {
			m := <-userChange
			userConf.Save(m)
		}
	}()

	// 4. setup a list of privileged users
	var userAllowed = make(map[string]string)
	userAllowed[config.Bot.Admin] = "all"

	// 5. setup service to retrive jokes
	jokeGetter := joke.NewJokeGetter()

	// 6. make connection to Telegram
	bot, err := tgbotapi.NewBotAPI(config.Bot.UUID)

	if config.Proxy.ProxyURL != "" {
		// Create Client to connect via Proxy
		logger.Info("Using proxy to connect")

		proxyURL := url.URL{
			Scheme: "socks5",
			User:   url.UserPassword(config.Proxy.User, config.Proxy.Password),
			Host:   config.Proxy.ProxyURL}

		transport := &http.Transport{
			Proxy: http.ProxyURL(&proxyURL),
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: transport}
		bot, err = tgbotapi.NewBotAPIWithClient(config.Bot.UUID, client)
	}

	if err != nil {
		logger.Panic("failed to connect to Telegram", err)
	}
	logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	go func() {
		for update := range updates {
			// универсальный ответ на любое сообщение
			reply := "Не знаю что сказать"
			if update.Message == nil {
				continue
			}

			user := update.Message.From.UserName

			// логируем от кого какое сообщение пришло
			logger.Info(fmt.Sprintf("[%s] %s", user, update.Message.Text))

			if update.Message.Text == "Шутка" {
				reply, err = jokeGetter.GetJoke()
				if err != nil {
					reply = "Error while getting a joke. Sorry.."
				}
			}

			if update.Message.Text == "Подписка" {
				reply = user + ", оформлена подписка на получение шуток каждые 2 часа"
				regUsers[user] = update.Message.Chat.ID
				userChange <- &regUsers
			}

			if update.Message.Text == "Отменить подписку" {
				reply = user + ", подписка отменена"
				delete(regUsers, user)
				userChange <- &regUsers
			}
			// In case the some file
			if update.Message.Text == "" {
				if update.Message.Document != nil {
					fileID := update.Message.Document.FileID
					fileName := update.Message.Document.FileName
					if !strings.Contains(fileName, ".torrent") {
						continue
					} else {
						if userAllowed[user] == "" {
							reply = "You are not authorized"
							// создаем ответное сообщение
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
							// отправляем
							bot.Send(msg)
							continue
						}

						url, err := bot.GetFileDirectURL(fileID)
						if err != nil {
							logger.Error("Failed to get file by fileID: "+fileID+" fileName: "+fileName, err)
						}

						response, err := http.Get(url)
						if err != nil {
							logger.Error("Error while downloading "+url, err)
							continue
						}
						defer response.Body.Close()

						data, err := ioutil.ReadAll(response.Body)
						if err != nil {
							logger.Error("Failed to get data from *Response.Body, err: %s", err)
						}
						var delugeConfig deluge.Deluge
						delugeConfig.Host = config.Deluge.Host
						delugeConfig.Port = config.Deluge.Port
						delugeConfig.Password = config.Deluge.Password
						err = deluge.AddTorrent(data, fileName, &delugeConfig)
						if err == nil {
							reply = "Torrent file " + fileName + " added"
						}
					}
				}
			}
			// свитч на обработку комманд
			// комманда - сообщение, начинающееся с "/"
			switch update.Message.Command() {
			case "start":
				reply = "Привет. Я телеграм-бот"
			}

			// создаем ответное сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

			if regUsers[user] != 0 {
				msg.ReplyMarkup = regUserKeyb
			} else {
				msg.ReplyMarkup = newUserKeyb
			}

			// отправляем
			bot.Send(msg)
		}
	}()

	// Рассылка погоды
	go func() {
		timer := make(chan struct{})
		go func() {
			time.Sleep(schedule.FirstStart(10))
			timer <- struct{}{}

		}()
		go func() {
			time.Sleep(schedule.FirstStart(17))
			timer <- struct{}{}
		}()
		for {
			_ = <-timer
			var wConfig weather.Weather
			wConfig.Key = config.Weather.Key
			response, err := weather.GetCurrentWeather(&wConfig)
			if err != nil {
				logger.Error("Error while getting weather", err)
				time.Sleep(4 * time.Hour)
				continue
			}
			message := "Cейчас: " + fmt.Sprintf("%.1f", response.Current.TempC) +
				" " + response.Current.Condition.Text +
				"\nОщущается как: " +
				fmt.Sprintf("%.1f", response.Current.FeelslikeC) +
				"\nВлажность: " +
				fmt.Sprintf("%d", response.Current.Humidity) +
				"\nДавление: " +
				fmt.Sprintf("%.1f", response.Current.PressureMb)
			for user := range regUsers {
				msg := tgbotapi.NewMessage(regUsers[user], message)
				bot.Send(msg)
			}
			go func() {
				time.Sleep(24 * time.Hour)
				timer <- struct{}{}
			}()
		}

	}()

	// Рассылка шуток
	for {
		message, err := jokeGetter.GetJoke()
		if err != nil {
			logger.Error("Error on getting joke", err)
			message = "Error while getting a joke. Sorry.."
		}
		for user := range regUsers {
			msg := tgbotapi.NewMessage(regUsers[user], message)
			bot.Send(msg)
		}
		time.Sleep(2 * time.Hour)
	}

}
