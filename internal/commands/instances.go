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

type InstancesCommand struct {
	*kingpin.CmdClause
	Service string
	Tail    bool

	monitoring domain.MonitoringService
}

func Instances(monitoring domain.MonitoringService) *InstancesCommand {
	return &InstancesCommand{
		monitoring: monitoring,
	}
}

func (c *InstancesCommand) Bind(cli *kingpin.Application) *InstancesCommand {
	c.CmdClause = cli.Command("list-instances", "Show the status of all running instances.")
	c.CmdClause.Flag("service", "Only show instances for a given service.").Short('s').StringVar(&c.Service)
	c.CmdClause.Flag("tail", "Tail periodic queries to list instances.").Short('t').BoolVar(&c.Tail)
	return c
}

func (c *InstancesCommand) Execute(ctx context.Context) error {
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

func (c *InstancesCommand) execute(ctx context.Context, w io.Writer) error {
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

		if c.Service == "" {
			servers = append(servers, server)
		} else if c.Service == server.Service {
			servers = append(servers, server)
		}
	}

	servers.SortByIP()
	c.printTable(servers, w)

	return nil
}

func (c *InstancesCommand) printTable(servers []*domain.Server, w io.Writer) {
	data := make([][]string, len(servers))
	for i, server := range servers {
		if server.IsHealthy() {
			data[i] = []string{
				server.IP.String(),
				server.Service,
				"Healthy",
				fmt.Sprintf("%.2f%%", server.CPUPercent),
				fmt.Sprintf("%.2f%%", server.MemoryPercent),
			}
		} else {
			data[i] = []string{
				color.RedString(server.IP.String()),
				color.RedString(server.Service),
				color.RedString("Unhealthy"),
				color.RedString(fmt.Sprintf("%.2f%%", server.CPUPercent)),
				color.RedString(fmt.Sprintf("%.2f%%", server.MemoryPercent)),
			}
		}
	}

	table := tablewriter.NewWriter(w)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"IP", "Service", "Status", "% CPU", "% Memory"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
