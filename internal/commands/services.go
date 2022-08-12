package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
)

type ServicesCommand struct {
	*kingpin.CmdClause
	Tail bool

	monitoring domain.MonitoringService
}

func Services(monitoring domain.MonitoringService) *ServicesCommand {
	return &ServicesCommand{
		monitoring: monitoring,
	}
}

func (c *ServicesCommand) Bind(cli *kingpin.Application) *ServicesCommand {
	c.CmdClause = cli.Command("list-services", "Show the status of all running services.")
	c.CmdClause.Flag("tail", "Tail periodic queries to list services.").Short('t').BoolVar(&c.Tail)
	return c
}

func (c *ServicesCommand) Execute(ctx context.Context) error {
	if !c.Tail {
		return c.execute(ctx, os.Stdout)
	}

	writer := uilive.New()
	writer.Start()
	w := writer.Newline()

	if err := c.execute(ctx, w); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			writer.Stop()
			return nil
		case <-time.After(time.Second * 5):
			err := c.execute(ctx, w)
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
		}
	}
}

func (c *ServicesCommand) execute(ctx context.Context, w io.Writer) error {
	ips, err := c.monitoring.ListServers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list servers: %s", err)
	}

	var servers domain.Servers
	for _, ip := range ips {
		server, err := c.monitoring.GetServer(ctx, ip)
		if err != nil {
			return fmt.Errorf("failed to get data about server %s: %s", ip, err)
		}
		servers = append(servers, server)
	}

	servers.SortByIP()
	services := servers.GroupByService()
	services.SortByName()
	c.printTable(services, w)

	return nil
}

func (c *ServicesCommand) printTable(services domain.Services, w io.Writer) {
	data := make([][]string, len(services))
	for i, service := range services {
		if service.NumHealthy < 2 {
			data[i] = []string{
				color.RedString(service.Name),
				color.RedString(fmt.Sprintf("%d/%d", service.NumHealthy, len(service.Instances))),
				color.RedString(fmt.Sprintf("%.2f%%", service.AverageCPUPercent)),
				color.RedString(fmt.Sprintf("%.2f%%", service.AverageMemoryPercent)),
			}
		} else {
			data[i] = []string{
				service.Name,
				fmt.Sprintf("%d/%d", service.NumHealthy, len(service.Instances)),
				fmt.Sprintf("%.2f%%", service.AverageCPUPercent),
				fmt.Sprintf("%.2f%%", service.AverageMemoryPercent),
			}
		}
	}

	table := tablewriter.NewWriter(w)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"Service", "Healthy", "% CPU Avg", "% Memory Avg"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
