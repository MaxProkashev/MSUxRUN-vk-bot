package main

import (
	"context"
	"log"
	"msuxrun-bot/internal/config"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

const packageErr = "ERROR cmd/bot/main.go::" // package standart error

type packErr struct {
	funcName string // name func where error
	comment  string // comment to error
	err      error  // err
}

func (p *packErr) addComment(comm string) {
	p.comment = comm
}

func errInit(pe packErr) {
	log.Fatalf("%s%s ==> %s <== COMMENT: %s",
		packageErr,
		pe.funcName,
		pe.err.Error(),
		pe.comment,
	)
}

var (
	stdErr packErr
	// standart output error for main func
	stdFuncErr = func(err error) packErr {
		return packErr{
			funcName: "MAIN",
			comment:  "no comment",
			err:      err,
		}
	}
)

func main() {
	var err error
	//? get configuration for bot
	conf := config.GetProjectConfig()
	conf.StartLog()

	//? create new vk api
	vk := api.NewVK(conf.Token)

	//? get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		stdErr = stdFuncErr(err)
		stdErr.addComment("can`t get info about group")
		errInit(stdErr)
	}
	for i, gr := range group {
		log.Printf("group[%d].id = %d", i, gr.ID)
	}

	//? initializing Long Poll in conf.GroupID
	lp, err := longpoll.NewLongPoll(vk, conf.GroupID)
	if err != nil {
		stdErr = stdFuncErr(err)
		stdErr.addComment("can`t get info about group")
		errInit(stdErr)
	}

	//? new message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

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
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		stdErr = stdFuncErr(err)
		stdErr.addComment("can`t run long poll")
		errInit(stdErr)
	}
}
