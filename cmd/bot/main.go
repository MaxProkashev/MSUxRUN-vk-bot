package main

import (
	"context"
	"log"
	"msuxrun-bot/internal/config"
	"msuxrun-bot/internal/logs"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func main() {
	//! start loggers
	logs.InitLoggers()

	//? get configuration for bot
	conf := config.GetProjectConfig()

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
		logs.Mess("from:%d text: %s", obj.Message.PeerID, obj.Message.Text)

		if obj.Message.Text == "ping" {
			b := params.NewMessagesSendBuilder()
			b.Message("pong")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)

			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	//? run bot Long Poll
	logs.Succes("start long poll with gr_id=%d, serv=%s, wait=%d", lp.GroupID, lp.Server, lp.Wait)
	if err := lp.Run(); err != nil {
		logs.ErrorLogger.Printf("could`t start long poll")
		os.Exit(1)
	}
}
