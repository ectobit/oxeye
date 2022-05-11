// Package service contains multithreaded worker pool implementation.
package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/broker"
	"go.ectobit.com/oxeye/encdec"
)

// Errors.
var (
	ErrInvalidMessageType = errors.New("invalid message type")
)

// Job defines common job methods.
type Job[IN, OUT any] interface {
	Execute(msg *IN) *OUT
}

// Service is a multithreaded service with configurable job to be executed.
type Service[IN, OUT any] struct {
	concurrency uint8
	broker      broker.Broker
	wg          sync.WaitGroup
	job         Job[IN, OUT]
	ed          encdec.EncDecoder
	log         lax.Logger
}

// NewService creates new service.
func NewService[IN, OUT any](concurrency uint8, broker broker.Broker, job Job[IN, OUT], endDec encdec.EncDecoder,
	log lax.Logger,
) *Service[IN, OUT] {
	return &Service[IN, OUT]{
		concurrency: concurrency,
		broker:      broker,
		job:         job,
		ed:          endDec,
		log:         log,
	}
}

// Run executes service reacting on termination signals for graceful shutdown.
func (s *Service[IN, OUT]) Run() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.log.Info("starting worker pool", lax.Uint8("concurrency", s.concurrency))

	sub, err := s.broker.Sub(ctx)
	if err != nil {
		return fmt.Errorf("broker: %w", err)
	}

	for workerID := uint8(1); workerID <= s.concurrency; workerID++ {
		go s.run(ctx, workerID, sub)
	}

	<-signals
	s.log.Info("graceful shutdown")
	s.wg.Wait()

	return nil
}

func (s *Service[IN, OUT]) run(ctx context.Context, workerID uint8, messages <-chan broker.Message) {
	s.log.Info("started", lax.Uint8("worker", workerID))
	s.wg.Add(1)

	for {
		select {
		case msg := <-messages:
			s.log.Debug("executing", lax.Uint8("worker", workerID))

			var inMsg IN

			if err := s.ed.Decode(msg.Data, &inMsg); err != nil {
				s.log.Warn("decode", lax.String("type", fmt.Sprintf("%T", inMsg)),
					lax.Uint8("worker", workerID), lax.Error(err))

				continue
			}

			msg.InProgress()

			outMsg := s.job.Execute(&inMsg)

			if outMsg == nil {
				continue
			}

			out, err := s.ed.Encode(outMsg)
			if err != nil {
				s.log.Warn("encode", lax.String("type", fmt.Sprintf("%T", outMsg)),
					lax.Uint8("worker", workerID), lax.Error(err))

				continue
			}

			if err := s.broker.Pub(out); err != nil {
				s.log.Warn("broker", lax.Uint8("worker", workerID), lax.Error(err))
			}

			msg.Ack()
		case <-ctx.Done():
			s.log.Info("stopped", lax.Uint8("worker", workerID))
			s.wg.Done()

			return
		}
	}
}

// Exit exits CLI application writing message and error to stderr.
func Exit(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
