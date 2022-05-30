// Package service contains multithreaded worker pool implementation.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	wg          sync.WaitGroup
	job         Job[IN, OUT]
	Debug       io.StringWriter
}

// NewService creates new service.
func NewService[IN, OUT any](concurrency uint8, broker broker.Broker, job Job[IN, OUT]) *Service[IN, OUT] {
	return &Service[IN, OUT]{
		concurrency: concurrency,
		broker:      broker,
		job:         job,
		Debug:       io.Discard.(io.StringWriter),
	}
}

// Run executes service reacting on termination signals for graceful shutdown.
func (s *Service[IN, OUT]) Run() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Debug.WriteString(fmt.Sprintf("starting worker pool with %d workers", s.concurrency))

	sub, err := s.broker.Sub(ctx)
	if err != nil {
		return fmt.Errorf("broker: %w", err)
	}

	for workerID := uint8(1); workerID <= s.concurrency; workerID++ {
		go s.run(ctx, workerID, sub)
	}

	<-signals
	s.Debug.WriteString("graceful shutdown")
	s.wg.Wait()

	return nil
}

func (s *Service[IN, OUT]) run(ctx context.Context, workerID uint8, messages <-chan broker.Message) {
	s.Debug.WriteString(fmt.Sprintf("starting worker %d", workerID))
	s.wg.Add(1)

	for {
		select {
		case msg := <-messages:
			s.Debug.WriteString(fmt.Sprintf("worker %d executing job", workerID))

			var inMsg IN

			if err := json.Unmarshal(msg.Data, &inMsg); err != nil {
				s.Debug.WriteString(fmt.Sprintf("worker %d decoding message type %T: %v", workerID, inMsg, err))

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
				s.Debug.WriteString(fmt.Sprintf("worker %d encoding message type %T: %v", workerID, outMsg, err))

				continue
			}

			if err := s.broker.Pub(out); err != nil {
				s.Debug.WriteString(fmt.Sprintf("worker %d publishing message %v: %v", workerID, inMsg, err))
			}

			msg.Ack()
		case <-ctx.Done():
			s.Debug.WriteString(fmt.Sprintf("stopping worker %d", workerID))
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
