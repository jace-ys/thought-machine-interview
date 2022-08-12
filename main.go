package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/thought-machine-interview/internal/commands"
	"github.com/jace-ys/thought-machine-interview/internal/httpapi"
)

func main() {
	cli := kingpin.New("cpxctl", "CLI tool for querying the health of services running in Cloud Provider X.")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cpx := httpapi.NewCPXClient()
	instancesCmd := commands.Instances(cpx).Bind(cli)
	servicesCmd := commands.Services(cpx).Bind(cli)

	go func() {
		<-ctx.Done()
		stop()
	}()

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case instancesCmd.FullCommand():
		cli.FatalIfError(instancesCmd.Execute(ctx), instancesCmd.FullCommand())
	case servicesCmd.FullCommand():
		cli.FatalIfError(servicesCmd.Execute(ctx), servicesCmd.FullCommand())
	}
}
