package main

import (
	"os"

	"github.com/nats-io/nats.go"
	"go.ectobit.com/act"
	"go.ectobit.com/lax"
	"go.ectobit.com/oxeye/broker"
	"go.ectobit.com/oxeye/encdec"
	"go.ectobit.com/oxeye/service"
)

type config struct {
	Env         string `help:"environment [dev|local|prod]" def:"local"`
	Concurrency uint   `def:"5"`
	NATS        struct {
		ClusterURL     act.URL `def:"nats://nats:4222"`
		ConsumeSubject string  `def:"OXEYE.in"`  // modify this to correct receive channel
		ConsumerGroup  string  `def:"oxeye"`     // modify this to correct consumer group
		ProduceSubject string  `def:"OXEYE.out"` // modify this to correct send channel
	}
	Log struct {
		Level  string `help:"log level [debug|info|warn|error]" def:"debug"`
		Format string `help:"console|json" def:"console"`
	}
	// add more configuration if needed
}

// This is en example service doing nothing useful. For the real microservice implementation you can clone
// this repository and use most of the files from the root. However, after cloning you should delete broker
// and service directories (packages) because they should be imported from this project, just like in the imports
// above. Also delete go.mod and go.sum and then run appropriate `go mod init`. So, you can keep useful files and
// directories like .vscode, .dockerignore, .drone.yml, .gitignore, docker-compose.yml, Dockerfile, job.go, main.go,
// Makefile and README.md. Of course, you should edit all of them and adapt for your microservice.
// Your job implementation should be in job.go and your main.go will probably be very similar to the provided one.
func main() {
	cfg := &config{} //nolint:exhaustivestruct

	// change service name to a proper one
	cli := act.New("oxeye")

	if err := cli.Parse(cfg, os.Args[1:]); err != nil {
		service.Exit("parsing flags", err)
	}

	nconn, err := nats.Connect(cfg.NATS.ClusterURL.String())
	if err != nil {
		service.Exit("nats connect", err)
	}

	jetStream, err := nconn.JetStream(nats.PublishAsyncMaxPending(256)) //nolint:gomnd
	if err != nil {
		service.Exit("nats jet stream", err)
	}

	// Create stream (check if this is idempotent)
	// _, err := js.AddStream(&nats.StreamConfig{
	// 	Name:     "OXEYE",
	// 	Subjects: []string{"OXEYE.*"},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	log, err := lax.NewDefaultZapAdapter(cfg.Log.Format, cfg.Log.Level)
	if err != nil {
		service.Exit("create logger", err)
	}

	brConfig := &broker.NatsJetStreamConfig{ //nolint:exhaustivestruct
		ConsumeSubject: cfg.NATS.ConsumeSubject,
		ConsumerGroup:  cfg.NATS.ConsumerGroup,
		ProduceSubject: cfg.NATS.ProduceSubject,
	}
	br := broker.NewNatsJetStream(jetStream, brConfig, log)

	srv := service.NewService[*InMsg, *OutMsg](uint8(cfg.Concurrency), br, &Job{log: log}, &encdec.JSONIter{}, log)

	if err := srv.Run(); err != nil {
		service.Exit("service run", err)
	}

	nconn.Close()

	log.Flush()
}
