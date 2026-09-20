package main

import (
	"context"
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
)

func TestRunEtlCallsTheInputPort(t *testing.T) {
	port := &importPortSpy{}
	if err := runEtl(context.Background(), port); err != nil {
		t.Fatalf("runEtl() error = %v", err)
	}
	if !port.called {
		t.Fatal("input port was not called")
	}
}

func TestRunEtlReturnsInputPortError(t *testing.T) {
	expected := errors.New("import failed")
	if err := runEtl(context.Background(), &importPortSpy{err: expected}); !errors.Is(err, expected) {
		t.Fatalf("runEtl() error = %v, want %v", err, expected)
	}
}

type importPortSpy struct {
	called bool
	err    error
}

func (s *importPortSpy) Execute(context.Context, importcatalog.Command) error {
	s.called = true
	return s.err
}
