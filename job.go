package main

import (
	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/service"
)

type InMsg struct{}

type OutMsg struct{}

var _ service.Job[InMsg, OutMsg] = (*Job)(nil)

type Job struct {
	log lax.Logger
}

func (j *Job) Execute(msg *InMsg) *OutMsg {
	// do something with msg
	_ = msg

	return &OutMsg{}
}
