package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ainvaltin/nu-plugin"
)

func main() {
	p, err := nu.New(
		[]*nu.Command{
			fromLogFmt(),
			toLogFmt(),
		},
		"0.1.1",
		nil,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create plugin", err)
		return
	}
	if err := p.Run(quitSignalContext()); err != nil && !errors.Is(err, nu.ErrGoodbye) {
		fmt.Fprintln(os.Stderr, "plugin exited with error", err)
	}
}

func quitSignalContext() context.Context {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigChan)
		sig := <-sigChan
		cancel(fmt.Errorf("got quit signal: %s", sig))
	}()

	return ctx
}
