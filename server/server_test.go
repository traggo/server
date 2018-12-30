package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/test"
)

func TestShutdownOnErrorWhileShutdown(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	disposeInterrupt := fakeInterrupt(t)
	defer disposeInterrupt()

	shutdownError := errors.New("shutdown error")
	disposeShutdown := fakeShutdownError(shutdownError)
	defer disposeShutdown()

	finished := make(chan error)

	go func() {
		finished <- Start(http.NewServeMux(), freeport.GetPort())
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Server should be closed")
	case err := <-finished:
		assert.Equal(t, shutdownError, err)
	}
}

func TestShutdownAfterError(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	finished := make(chan error)

	go func() {
		finished <- Start(http.NewServeMux(), -5)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Server should be closed")
	case err := <-finished:
		assert.NotNil(t, err)
	}
}

func TestShutdown(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	dispose := fakeInterrupt(t)
	defer dispose()

	finished := make(chan error)

	go func() {
		finished <- Start(http.NewServeMux(), freeport.GetPort())
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Server should be closed")
	case err := <-finished:
		assert.Nil(t, err)
	}
}

func fakeInterrupt(t *testing.T) func() {
	oldNotify := notifySignal
	notifySignal = func(c chan<- os.Signal, sig ...os.Signal) {
		assert.Contains(t, sig, os.Interrupt)
		go func() {
			select {
			case <-time.After(100 * time.Millisecond):
				c <- os.Interrupt
			}
		}()
	}
	return func() {
		notifySignal = oldNotify
	}
}

func fakeShutdownError(err error) func() {
	old := serverShutdown
	serverShutdown = func(server *http.Server, ctx context.Context) error {
		return err
	}
	return func() {
		serverShutdown = old
	}
}
