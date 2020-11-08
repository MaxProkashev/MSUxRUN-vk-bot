package logbot

import "log"

const (
	packageMainErr   = "ERROR cmd/bot/main.go::"           // package main standart error
	packageConfigErr = "ERROR internal/config/config.go::" // package config standart error
)

var StdFuncErr = func(err error) PackErr {
	return PackErr{
		FuncName: "GetProjectConfig",
		Comment:  "no comment",
		Err:      err,
	}
}

type PackErr struct {
	PackName string
	FuncName string
	Comment  string
	Err      error
}

func (p *PackErr) AddComment(comm string) {
	p.Comment = comm
}

func (p *PackErr) ERRinit() {
	log.Fatalf("%s%s ==> %s <== COMMENT: %s",
		p.PackName,
		p.FuncName,
		p.Err.Error(),
		p.Comment,
	)
}
