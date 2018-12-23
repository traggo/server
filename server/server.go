package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/gql"
)

var notifySignal = signal.Notify
var serverShutdown = func(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}

// Start starts the http server
func Start(db *gorm.DB, passwordStrength, port int) error {
	server, shutdown := startServer(db, passwordStrength, port)
	shutdownOnInterruptSignal(server, 2*time.Second, shutdown)
	return waitForServerToClose(shutdown)
}

func startServer(db *gorm.DB, passwordStrength, port int) (*http.Server, chan error) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.Handle("/graphql", gql.Handler(db, passwordStrength))

	shutdown := make(chan error)

	go func() {
		err := srv.ListenAndServe()
		shutdown <- err
	}()
	return srv, shutdown
}

func shutdownOnInterruptSignal(server *http.Server, timeout time.Duration, shutdown chan<- error) {
	interrupt := make(chan os.Signal, 1)
	notifySignal(interrupt, os.Interrupt)

	go func() {
		select {
		case <-interrupt:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := serverShutdown(server, ctx); err != nil {
				shutdown <- err
			}
		}
	}()
}

func waitForServerToClose(shutdown <-chan error) error {
	err := <-shutdown
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
