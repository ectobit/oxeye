package main

import (
	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/service"
)

// InMsg - replace it with your input message.
type InMsg struct{}

// OutMsg - replace it with your output message.
type OutMsg struct{}

var _ service.Job[InMsg, OutMsg] = (*Job)(nil)

// Job contains job dependencies.
type Job struct {
	log lax.Logger
}

// Execute executes job.
func (j *Job) Execute(msg *InMsg) *OutMsg {
	// do something with msg
	_ = msg

	return &OutMsg{}
}
