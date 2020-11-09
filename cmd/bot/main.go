package main

import (
	"context"
	"msuxrun-bot/internal/config"
	"msuxrun-bot/internal/db"
	"msuxrun-bot/internal/logs"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
)

var val int

func main() {
	//! start loggers
	logs.InitLoggers()

	//! get configuration for bot and db
	conf := config.GetProjectConfig()
	db.DB = conf.OpenDB().DB

	//? create bot_users table
	db.DropTable("bot_users")
	db.CreateUserTable()

	//? create new vk api
	vk := api.NewVK(conf.Token)

	//? get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		logs.ErrorLogger.Printf("can`t get info about group")
		os.Exit(1)
	}
	for i, gr := range group {
		logs.Succes("group[%d].id = %d", i, gr.ID)
	}

	//? initializing Long Poll in conf.GroupID
	lp, err := longpoll.NewLongPoll(vk, conf.GroupID)
	if err != nil {
		logs.ErrorLogger.Printf("could`t init long poll in group_id = %d", conf.GroupID)
		os.Exit(1)
	}

	//? new message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		logs.Mess("from: %d, text: %s, payload: %s",
			obj.Message.FromID,
			obj.Message.Text,
			obj.Message.Payload,
		)
		//? get info about user
		user := &db.User{
			ID: obj.Message.FromID,
		}
		user = user.CheckUserByID()

		stdKEY := createUserKey(user, conf)
		b := params.NewMessagesSendBuilder()
		switch obj.Message.Text {
		case "Лонгран ПН 19:00":
			val = db.GetInt(user.ID, "mo")
			if val == 0 {
				b.Message("Ты записан на тренировку в понедельник в 19:00")
				db.SetInt(user.ID, "mo", 1)
				user.MO = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "mo", 0)
				user.MO = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Общая ВТ 19:30":
			val = db.GetInt(user.ID, "tu")
			if val == 0 {
				b.Message("Ты записан на тренировку во вторник в 19:30")
				db.SetInt(user.ID, "tu", 1)
				user.TU = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "tu", 0)
				user.TU = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Лонгран СР 19:00":
			val = db.GetInt(user.ID, "we")
			if val == 0 {
				b.Message("Ты записан на тренировку в среду в 19:00")
				db.SetInt(user.ID, "we", 1)
				user.WE = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "we", 0)
				user.WE = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Общая ЧТ 19:30":
			val = db.GetInt(user.ID, "th")
			if val == 0 {
				b.Message("Ты записан на тренировку в четверг в 19:30")
				db.SetInt(user.ID, "th", 1)
				user.TH = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "th", 0)
				user.TH = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Лонгран ПТ 19:00":
			val = db.GetInt(user.ID, "fr")
			if val == 0 {
				b.Message("Ты записан на тренировку в пятницу в 19:00")
				db.SetInt(user.ID, "fr", 1)
				user.FR = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "fr", 0)
				user.FR = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Лонгран ВС 9:00(9:20)":
			val = db.GetInt(user.ID, "su")
			if val == 0 {
				b.Message("Ты записан на тренировку в воскресенье в 19:00")
				db.SetInt(user.ID, "su", 1)
				user.SU = 1
			} else {
				b.Message("Ты отписался от тренировки")
				db.SetInt(user.ID, "su", 0)
				user.SU = 0
			}
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(createUserKey(user, conf).ToJSON())
		case "Расписание":
			b.Message("Вот расписание на неделю")
			b.Attachment("photo-186543814_457239215")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(stdKEY.ToJSON())
		case "Мои тренировки":
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(stdKEY.ToJSON())
			str := ""
			if user.MO == 1 {
				str += "ПН 19:00\nЛонгран\nТренеры: Сергей Павлов, Андрей Королев\n\n"
			}
			if user.TU == 1 {
				str += "ВТ 19:30\nОбщая тренировка\nТренер: Артем Орлов\n\n"
			}
			if user.WE == 1 {
				str += "СР 19:00\nЛонгран\nТренеры: Сергей Павлов, Андрей Королев\n\n"
			}
			if user.TH == 1 {
				str += "ЧТ 19:30\nОбщая тренировка\nТренер: Артём Орлов\n\n"
			}
			if user.FR == 1 {
				str += "ПТ 19:00\nЛонгран\nТренеры: Сергей Павлов, Андрей Королев\n\n"
			}
			if user.SU == 1 {
				str += "ВС 9:00(9:20)\nЛонгран\nТренеры: Сухайли Хотамов\n\n"
			}
			if str == "" {
				str += "У тебя еще нет ни одной тренировки"
			}
			b.Message(str)
		default:
			//? create default message to user
			b.Message("Лучше пообщаемся на тренировке!")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			b.Keyboard(stdKEY.ToJSON())
		}

		_, err := vk.MessagesSend(b.Params)
		if err != nil {
			logs.Warn("can`t send message. Reason: %s", err.Error())
		}

	})

	//? run bot Long Poll
	logs.Succes("start long poll with gr_id=%d, serv=%s, wait=%d", lp.GroupID, lp.Server, lp.Wait)
	if err := lp.Run(); err != nil {
		logs.ErrorLogger.Printf("could`t start long poll")
		os.Exit(1)
	}
}

func createUserKey(user *db.User, conf config.Config) (stdKEY *object.MessagesKeyboard) {
	//? create standart hotkey
	stdKEY = object.NewMessagesKeyboard(false)

	stdKEY.AddRow()

	if user.MO == 1 {
		stdKEY.AddTextButton(
			"Лонгран ПН 19:00",
			"mo",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Лонгран ПН 19:00",
			"mo",
			"secondary",
		)
	}

	if user.TU == 1 {
		stdKEY.AddTextButton(
			"Общая ВТ 19:30",
			"tu",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Общая ВТ 19:30",
			"tu",
			"secondary",
		)
	}

	stdKEY.AddRow()

	if user.WE == 1 {
		stdKEY.AddTextButton(
			"Лонгран СР 19:00",
			"we",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Лонгран СР 19:00",
			"we",
			"secondary",
		)
	}

	if user.TH == 1 {
		stdKEY.AddTextButton(
			"Общая ЧТ 19:30",
			"th",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Общая ЧТ 19:30",
			"th",
			"secondary",
		)
	}

	stdKEY.AddRow()

	if user.FR == 1 {
		stdKEY.AddTextButton(
			"Лонгран ПТ 19:00",
			"fr",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Лонгран ПТ 19:00",
			"fr",
			"secondary",
		)
	}

	if user.SU == 1 {
		stdKEY.AddTextButton(
			"Лонгран ВС 9:00(9:20)",
			"su",
			"positive",
		)
	} else {
		stdKEY.AddTextButton(
			"Лонгран ВС 9:00(9:20)",
			"su",
			"secondary",
		)
	}

	stdKEY.AddRow()
	stdKEY.AddTextButton(
		"Расписание",
		"Расписание",
		"primary",
	)

	stdKEY.AddRow()
	stdKEY.AddTextButton(
		"Мои тренировки",
		"Мои тренировки",
		"primary",
	)
	return stdKEY
}
