package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alexey-s-aksenov/tbot/schedule"
	"github.com/alexey-s-aksenov/tbot/weather"

	"github.com/alexey-s-aksenov/tbot/saveusers"

	"github.com/alexey-s-aksenov/tbot/deluge"
	"github.com/alexey-s-aksenov/tbot/joke"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"gopkg.in/yaml.v2"
)

var config tbotConfig

type tbotConfig struct {
	Bot struct {
		UUID  string `yaml:"uuid"`
		Admin string `yaml:"admin"`
	}
	Log struct {
		File string `yaml:"file"`
	}
	Deluge struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	}
	Weather struct {
		Key string `yaml:"key"`
	}
	Users struct {
		File string `yaml:"file"`
	}
	Proxy struct {
		ProxyURL string `yaml:"proxy"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
}

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

var userAllowed = make(map[string]string)
var regUsers = make(map[string]int64)
var userConf saveusers.EncConfig

func init() {

	f, err := os.OpenFile("tbot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	file, err := os.Open("config.yaml")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file)

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Printing config:\nAdmin: %s", config.Bot.Admin)
	userAllowed[config.Bot.Admin] = "all"

	userConf.File = config.Users.File
	users, err := userConf.Load()

	if err != nil {
		log.Panic(err)
	}

	regUsers = *users
}

func main() {

	f, err := os.OpenFile("tbot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	jokeGetter := joke.NewJokeGetter()

	bot, err := tgbotapi.NewBotAPI(config.Bot.UUID)

	if config.Proxy.ProxyURL != "" {
		// Create Client to connect via Proxy
		log.Println("Using proxy to connect")

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
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	userChange := make(chan *map[string]int64)
	go func() {
		for {
			m := <-userChange
			userConf.Save(m)
		}
	}()

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
			log.Printf("[%s] %s", user, update.Message.Text)

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
							log.Printf("Failed to get file by fileID: " + fileID + " fileName: " + fileName)
						}

						response, err := http.Get(url)
						if err != nil {
							log.Printf("Error while downloading %s - %s", url, err)
							continue
						}
						defer response.Body.Close()

						data, err := ioutil.ReadAll(response.Body)
						if err != nil {
							log.Printf("Failed to get data from *Response.Body, err: %s", err)
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
				log.Printf("Error while getting weather %s", err)
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
			message = "Error while getting a joke. Sorry.."
		}
		for user := range regUsers {
			msg := tgbotapi.NewMessage(regUsers[user], message)
			bot.Send(msg)
		}
		time.Sleep(2 * time.Hour)
	}
}
