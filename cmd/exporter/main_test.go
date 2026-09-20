package main

import (
	"context"
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
)

func TestRunExporterInvokesTheInputPort(t *testing.T) {
	port := &exportPortSpy{}

	if err := runExporter(context.Background(), port); err != nil {
		t.Fatalf("runExporter() error = %v", err)
	}
	if !port.called {
		t.Fatal("input port was not called")
	}
}

func TestRunExporterReturnsTheInputPortError(t *testing.T) {
	expected := errors.New("export failed")

	err := runExporter(context.Background(), &exportPortSpy{err: expected})

	if !errors.Is(err, expected) {
		t.Fatalf("runExporter() error = %v, want %v", err, expected)
	}
}

type exportPortSpy struct {
	called bool
	err    error
}

func (s *exportPortSpy) Execute(context.Context, exportavailability.Command) error {
	s.called = true
	return s.err
}
