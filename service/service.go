// Package service contains multithreaded worker pool implementation.
package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/broker"
	"go.ectobit.com/oxeye/encdec"
)

// Job defines common job methods.
type Job interface {
	Execute(msg interface{}, ack func()) interface{}
	NewInMessage() interface{}
}

// Service is a multithreaded service with configurable job to be executed.
type Service struct {
	concurrency uint8
	broker      broker.Broker
	wg          sync.WaitGroup
	job         Job
	ed          encdec.EncDecoder
	log         lax.Logger
}

// NewService creates new service.
func NewService(concurrency uint8, broker broker.Broker, job Job, ed encdec.EncDecoder, log lax.Logger) *Service {
	return &Service{ //nolint:exhaustivestruct
		concurrency: concurrency,
		broker:      broker,
		job:         job,
		ed:          ed,
		log:         log,
	}
}

// Run executes service reacting on termination signals for graceful shutdown.
func (s *Service) Run() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.log.Info("starting worker pool", lax.Uint8("concurrency", s.concurrency))

	sub, err := s.broker.Sub(ctx)
	if err != nil {
		return fmt.Errorf("broker: %w", err)
	}

	s.wg.Add(int(s.concurrency))

	for workerID := uint8(1); workerID <= s.concurrency; workerID++ {
		go func(ctx context.Context, workerID uint8) {
			s.log.Info("started", lax.Uint8("worker", workerID))

			for {
				select {
				case msg := <-sub:
					s.log.Debug("executing", lax.Uint8("worker", workerID))

					inMsg := s.job.NewInMessage()

					if err := s.ed.Decode(msg.Data, inMsg); err != nil {
						s.log.Warn("decode", lax.String("type", fmt.Sprintf("%T", inMsg)),
							lax.Uint8("worker", workerID), lax.Error(err))

						continue
					}

					outMsg := s.job.Execute(inMsg, msg.Ack)

					out, err := s.ed.Encode(outMsg)
					if err != nil {
						s.log.Warn("encode", lax.String("type", fmt.Sprintf("%T", outMsg)),
							lax.Uint8("worker", workerID), lax.Error(err))

						continue
					}

					if err := s.broker.Pub(out); err != nil {
						s.log.Warn("broker", lax.Uint8("worker", workerID), lax.Error(err))
					}
				case <-ctx.Done():
					s.log.Info("stopped", lax.Uint8("worker", workerID))
					s.wg.Done()

					return
				}
			}
		}(ctx, workerID)
	}

	<-signals
	s.log.Info("graceful shutdown")

	return nil
}

// Exit exits CLI application writing message and error to stderr.
func Exit(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v", message, err)
	os.Exit(1)
}
