package main

import (
	"context"
	"os"
	"time"

	ndt7 "github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
)

const (
	clientName     = "ndt7-client-go-cmd"
	clientVersion  = "0.5.0"
	defaultTimeout = 55 * time.Second
)

type startFunc func(context.Context) (<-chan spec.Measurement, error)

func main() {
	client := ndt7.NewClient(clientName, clientVersion)
	parentCtx := context.Background()

	ctx, cancelFunc := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancelFunc()

	tests := map[spec.TestKind]startFunc{
		spec.TestDownload: client.StartDownload,
		spec.TestUpload:   client.StartUpload,
	}

	e := NewEmitterOutput(os.Stdout)
	for spec, f := range tests {
		e.testRunner(ctx, spec, f)
	}

	e.Summary(client)
}

func (e EmitterOutput) testRunner(ctx context.Context, kind spec.TestKind, start startFunc) {
	ch, err := start(ctx)

	err = e.Started(kind)
	if err != nil {
		e.Failed(kind, err)
		os.Exit(1)
	}

	for event := range ch {
		func(m *spec.Measurement) {
			e.SpeedEvent(&event)
		}(&event)
	}

	err = e.Completed(kind)
	if err != nil {
		e.Failed(kind, err)
		os.Exit(1)
	}
}
