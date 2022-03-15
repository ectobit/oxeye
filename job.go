package main

import (
	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/broker"
	"go.ectobit.com/oxeye/encdec"
	"go.ectobit.com/oxeye/service"
)

type InMsg struct{}

type OutMsg struct{}

var _ service.Job[*InMsg, *OutMsg] = (*Job)(nil)

type Job struct {
	broker broker.Broker
	ed     encdec.EncDecoder
	log    lax.Logger
}

func (j *Job) Execute(msg *InMsg) error {
	// do something with msg
	_ = msg

	// create output message from your result if needed
	outMsg := msg

	// encode your outMsg
	out, err := j.ed.Encode(outMsg)
	if err != nil {
		return err
	}

	// publish your outMsg
	if err := j.broker.Pub(out); err != nil {
		return err
	}

	return nil
}
