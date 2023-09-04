// Package service contains multithreaded worker pool implementation.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.ectobit.com/oxeye/broker"
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
	done        chan struct{}
	wg          sync.WaitGroup
	job         Job[IN, OUT]
	Debug       func(s string)
}

// NewService creates new service.
func NewService[IN, OUT any](concurrency uint8, broker broker.Broker, job Job[IN, OUT]) *Service[IN, OUT] {
	return &Service[IN, OUT]{ //nolint:exhaustruct
		concurrency: concurrency,
		broker:      broker,
		done:        make(chan struct{}),
		job:         job,
		Debug:       func(string) {},
	}
}

// Run executes service reacting on termination signals for graceful shutdown.
func (s *Service[IN, OUT]) Run() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	s.Debug(fmt.Sprintf("starting worker pool with %d workers", s.concurrency))

	sub, err := s.broker.Sub()
	if err != nil {
		return fmt.Errorf("broker: %w", err)
	}

	for workerID := uint8(1); workerID <= s.concurrency; workerID++ {
		go s.run(workerID, sub)
	}

	<-signals
	s.Debug("graceful shutdown")
	s.wg.Wait()
	s.broker.Exit()

	return nil
}

func (s *Service[IN, OUT]) run(workerID uint8, messages <-chan broker.Message) {
	s.Debug(fmt.Sprintf("starting worker %d", workerID))
	s.wg.Add(1)

	for {
		select {
		case msg := <-messages:
			s.Debug(fmt.Sprintf("worker %d executing job", workerID))

			var inMsg IN

			if err := json.Unmarshal(msg.Data, &inMsg); err != nil {
				s.Debug(fmt.Sprintf("worker %d decoding message type %T: %v", workerID, inMsg, err))

				continue
			}

			msg.InProgress()

			outMsg := s.job.Execute(&inMsg)

			if outMsg == nil {
				msg.Ack()

				continue
			}

			out, err := json.Marshal(outMsg)
			if err != nil {
				s.Debug(fmt.Sprintf("worker %d encoding message type %T: %v", workerID, outMsg, err))

				continue
			}

			if err := s.broker.Pub(out); err != nil {
				s.Debug(fmt.Sprintf("worker %d publishing message %v: %v", workerID, inMsg, err))
			}

			msg.Ack()
		case <-s.done:
			s.Debug(fmt.Sprintf("stopping worker %d", workerID))
			s.wg.Done()

			for range messages {
				<-messages
			}

			return
		}
	}
}

// Exit exits CLI application writing message and error to stderr.
func Exit(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
