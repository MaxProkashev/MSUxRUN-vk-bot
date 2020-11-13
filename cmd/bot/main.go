package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/MaxProkashev/MSUxRUN-vk-bot/internal/config"
	"github.com/MaxProkashev/MSUxRUN-vk-bot/internal/db"
	"github.com/MaxProkashev/MSUxRUN-vk-bot/internal/logs"
	"github.com/gin-gonic/gin"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
)

var (
	// only re-create table
	create = func() {
		db.CreateUserTable()
	}
	// first drop then create table
	createDrop = func() {
		db.DropUserTable()
		db.CreateUserTable()
	}
	// time for sleep when first call demon
	durFirstSleepy = func(h, m, s int) time.Duration {
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, msk)
		if t.Before(now) {
			t = t.Add(time.Hour * 24)
		}
		return t.Sub(time.Now())
	}

	wg            = &sync.WaitGroup{}
	mu            = &sync.Mutex{}
	msk           = time.FixedZone("UTC+3", +3*60*60)
	vk            *api.VK
	conf          *config.Config
	goroutinesNum = 8
	appURL        = "https://msu-vk-bot.herokuapp.com/"
)

func main() {
	runtime.GOMAXPROCS(0)

	// get port
	port := os.Getenv("PORT")
	if port == "" {
		logs.Err("%s", "without port")
		os.Exit(5)
	}
	// gin router for heroku

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	err := router.Run(":" + port)
	if err != nil {
		logs.Err("Could not run router. Reason: %s", err.Error())
	}

	//! start loggers
	logs.InitLoggers()
	//! get configuration for project
	conf = config.GetProjectConfig()
	//! init DB for bot
	db.InitDB(conf.DbURL)

	//? create bot_users table
	create()

	//! demon
	callAt(conf.CallH, conf.CallM, conf.CallS)

	// ? create new vk api
	vk = api.NewVK(conf.Token)

	// ? get information about the group
	groups, err := vk.GroupsGetByID(nil)
	if err != nil {
		logs.Err("can`t get info about group")
		os.Exit(1)
	}
	for _, gr := range groups {
		logs.Succes("group[%d] %s", gr.ID, gr.Name)

		wg.Add(1)
		go func(gr object.GroupsGroup, waiter *sync.WaitGroup) {
			defer waiter.Done()
			//? initializing Long Poll in conf.GroupID
			lp, err := longpoll.NewLongPoll(vk, gr.ID)
			if err != nil {
				logs.Err("could`t init long poll in %d. Reason: %s", gr.ID, err.Error())
				os.Exit(1)
			}

			//? reg new message event
			lp.MessageNew(standartMessageEvent)

			//? run n bot Long Poll
			logs.Succes("group[%d] start long poll", lp.GroupID)
			if err := lp.Run(); err != nil {
				logs.Err("could`t start long poll in %d. Reason: %s", gr.ID, err.Error())
				os.Exit(1)
			}
		}(gr, wg)
	}

	wg.Wait()
}

func callAt(h, m, s int) {
	duration := durFirstSleepy(h, m, s)
	//fmt.Println(duration)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		time.Sleep(duration) // засыпаем до первого вызова
		for {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				mu.Lock()
				postNot()
				mu.Unlock()
			}(wg)
			time.Sleep(time.Hour * 24) // засыпекм до следующего вызова
		}
	}(wg)
}

func postNot() {
	var msg string
	weekday := time.Now().Weekday().String()
	logs.WarnNot(time.Now().Format(time.RFC3339))
	rand.Seed(time.Now().UnixNano())

	var randNot = func() string {
		return conf.MessageNotice[rand.Intn(len(conf.MessageNotice))]
	}

	ch := make(chan *db.User, 1)
	go db.GetAllUser(ch)

	for user := range ch {
		user.ParseSign(conf.CountTrain)

	LOOP:
		for i, tr := range user.Train {
			if tr == true && conf.MainKeyboard[i].NotDay == weekday {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.PeerID(user.ID)
				msg = fmt.Sprintf(randNot(), conf.MainKeyboard[i].Label)
				b.Message(msg)

				_, err := vk.MessagesSend(b.Params)
				if err != nil {
					logs.Warn("can`t send message. Reason: %s", err.Error())
				}
				break LOOP
			}
		}
	}
}

// StandartMessageEvent from user
var standartMessageEvent = func(_ context.Context, obj events.MessageNewObject) {
	var msg string
	logs.Mess("user[%d] %s", obj.Message.FromID, obj.Message.Text)
	user := &db.User{
		Text: obj.Message.Text,
	}

	user.GetUser(obj.Message.FromID) // id,sign
	user.ParseSign(conf.CountTrain)  // train: [true false true false false false]

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(obj.Message.PeerID)
	// проверка кнопок записи на тренировку
	for i, fl := range conf.MainKeyboard {
		if user.Text == fl.Label {
			user.Train[i] = !user.Train[i]
			user.SetTrain(conf.CountTrain)
			if user.Train[i] {
				msg = conf.MessageSignUp + conf.MainKeyboard[i].Label
				if i == 0 {
					msg += conf.MessageSpecFor0
				}
			} else {
				msg = conf.MessageSignOut
			}
			b.Message(msg)
			b.Keyboard(renderKey(user.Train).ToJSON())

			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				logs.Warn("can`t send message. Reason: %s", err.Error())
			}
			return
		}
	}
	b.Keyboard(renderKey(user.Train).ToJSON())

	// проверка остальных действий
	switch user.Text {
	case conf.Schedule:
		msg = conf.Schedule
		b.Attachment(conf.SchPhoto)
	case conf.MyTrain:
		msg = userTrain(user.Train)
	default:
		msg = conf.MessageDefault
	}

	b.Message(msg)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		logs.Warn("can`t send message. Reason: %s", err.Error())
	}
	return
}

func renderKey(tr []bool) *object.MessagesKeyboard {
	var color string

	key := object.NewMessagesKeyboard(true)
	for i, fl := range tr {
		if i%2 == 0 {
			key.AddRow()
		}
		switch fl {
		case true:
			color = object.ButtonGreen
		case false:
			color = object.ButtonWhite
		}
		key.AddTextButton(
			conf.MainKeyboard[i].Label,
			"",
			color,
		)
	}

	key.AddRow()
	key.AddTextButton(
		conf.Schedule,
		"",
		object.ButtonBlue,
	)
	key.AddRow()
	key.AddTextButton(
		conf.MyTrain,
		"",
		object.ButtonBlue,
	)

	return key
}

func userTrain(tr []bool) string {
	str := ""

	for i, fl := range tr {
		if fl {
			str += fmt.Sprintf("%s%s",
				conf.MainKeyboard[i].Label,
				conf.MainKeyboard[i].Coach,
			)
		}
	}

	if str == "" {
		return conf.MessageNonTrain
	}
	return str
}
