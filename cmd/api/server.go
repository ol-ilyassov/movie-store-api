package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve
func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// Receives any errors returned by the graceful Shutdown() function.
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.logger.Info("completing background tasks", "addr", srv.Addr)

		app.wg.Wait()        // wait until waitGroup counter is zero.
		shutdownError <- nil // send that shutdown completed without any issues.
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}

// Signals:
// - SIGINT  - graceful shutdown     - signal: interrupt.
//             (keyboard)
// - SIGQUIT - exit without graceful shutdown  - signal: quit.
//             (keyboard) (with stack dump)
// - SIGTERM - graceful shutdown     - signal: terminated.
//             (terminate process in orderly manner)
// - SIGKILL - exit without graceful shutdown  - signal: killed.
//           - (uncatchable)
