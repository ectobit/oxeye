package main

import (
	"fmt"

	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/service"
)

type inMsg struct{}

type outMsg struct{}

var _ service.Job = (*job)(nil)

type job struct {
	log lax.Logger
}

func newJob(log lax.Logger) *job {
	return &job{log: log}
}

func (j *job) Execute(msg interface{}) (interface{}, error) {
	msg, ok := msg.(*inMsg)
	if !ok {
		return nil, fmt.Errorf("%w: %T", service.ErrInvalidMessageType, msg)
	}

	// do something with in
	_ = msg

	return &outMsg{}, nil
}

func (j *job) NewInMessage() interface{} {
	return &inMsg{}
}
