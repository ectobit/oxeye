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

func (j *job) Execute(msg interface{}, ack func()) interface{} {
	msg, ok := msg.(*inMsg)
	if !ok {
		j.log.Warn("assure message", lax.String("invalid type", fmt.Sprintf("%T", msg)))

		return nil
	}

	// do something with in
	_ = msg

	ack() // don't forget to call ack on success

	return &outMsg{}
}

func (j *job) NewInMessage() interface{} {
	return &inMsg{}
}

func (j *job) NewOutMessage() interface{} {
	return &outMsg{}
}
