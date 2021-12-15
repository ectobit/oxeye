package main

import (
	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/service"
)

type InMsg struct{}

type OutMsg struct{}

var _ service.Job[*InMsg, *OutMsg] = (*job)(nil)

type job struct {
	log lax.Logger
}

func (j *job) Execute(msg *InMsg) *OutMsg {
	// do something with msg
	_ = msg

	return &OutMsg{}
}
