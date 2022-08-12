package domain

import (
	"bytes"
	"context"
	"errors"
	"net"
	"sort"
)

var (
	ErrInvalidServerIP = errors.New("invalid IP address for server")
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MonitoringService
type MonitoringService interface {
	ListServers(ctx context.Context) ([]net.IP, error)
	GetServer(ctx context.Context, ip net.IP) (*Server, error)
}

type Servers []*Server

func (s Servers) SortByIP() {
	sort.Stable(s)
}

func (s Servers) Len() int           { return len(s) }
func (s Servers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Servers) Less(i, j int) bool { return bytes.Compare(s[i].IP, s[j].IP) < 0 }

type Server struct {
	IP            net.IP
	Service       string
	CPUPercent    float64
	MemoryPercent float64
}

func (s *Server) IsHealthy() bool {
	// We consider a service unhealthy if its CPU or memory utilisation is more than 90%
	switch {
	case s.CPUPercent > 90:
		return false
	case s.MemoryPercent > 90:
		return false
	default:
		return true
	}
}

func (s Servers) GroupByService() Services {
	sm := make(map[string]Servers)
	for _, server := range s {
		sm[server.Service] = append(sm[server.Service], server)
	}

	var services Services
	for name, instances := range sm {
		var healthy int
		var avgCPU float64
		var avgMem float64

		for _, instance := range instances {
			if instance.IsHealthy() {
				healthy++
			}

			avgCPU += instance.CPUPercent / float64(len(instances))
			avgMem += instance.MemoryPercent / float64(len(instances))
		}

		services = append(services, &Service{
			Name:                 name,
			Instances:            instances,
			NumHealthy:           healthy,
			AverageCPUPercent:    avgCPU,
			AverageMemoryPercent: avgMem,
		})
	}

	return services
}

type Services []*Service

func (s Services) SortByName() {
	sort.Stable(s)
}

func (s Services) Len() int           { return len(s) }
func (s Services) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Services) Less(i, j int) bool { return s[i].Name < s[j].Name }

type Service struct {
	Name                 string
	Instances            Servers
	NumHealthy           int
	AverageCPUPercent    float64
	AverageMemoryPercent float64
}
