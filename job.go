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
		j.log.Warn("assure message", lax.String("invalid type", fmt.Sprintf("%T", msg)))

		return nil, fmt.Errorf("%w: %T", service.ErrInvalidMessageType, msg)
	}

	// do something with in
	_ = msg

	return &outMsg{}, nil
}

func (j *job) NewInMessage() interface{} {
	return &inMsg{}
}

func (j *job) NewOutMessage() interface{} {
	return &outMsg{}
}
