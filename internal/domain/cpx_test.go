package domain_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
)

func TestGroupByService(t *testing.T) {
	servers := domain.Servers{
		{
			IP:            net.ParseIP("10.58.1.1"),
			Service:       "UserService",
			CPUPercent:    25.00,
			MemoryPercent: 75.00,
		},
		{
			IP:            net.ParseIP("10.58.1.2"),
			Service:       "GeoService",
			CPUPercent:    50.00,
			MemoryPercent: 50.00,
		},
		{
			IP:            net.ParseIP("10.58.1.3"),
			Service:       "StorageService",
			CPUPercent:    95.00,
			MemoryPercent: 65.00,
		},
		{
			IP:            net.ParseIP("10.58.1.4"),
			Service:       "UserService",
			CPUPercent:    10.00,
			MemoryPercent: 35.00,
		},
		{
			IP:            net.ParseIP("10.58.1.5"),
			Service:       "GeoService",
			CPUPercent:    40.00,
			MemoryPercent: 100.00,
		},
	}

	services := servers.GroupByService()

	assert.ElementsMatch(t, services, domain.Services{
		{
			Name:                 "UserService",
			Instances:            domain.Servers{servers[0], servers[3]},
			NumHealthy:           2,
			AverageCPUPercent:    17.50,
			AverageMemoryPercent: 55.00,
		},
		{
			Name:                 "GeoService",
			Instances:            domain.Servers{servers[1], servers[4]},
			NumHealthy:           1,
			AverageCPUPercent:    45.00,
			AverageMemoryPercent: 75.00,
		},
		{
			Name:                 "StorageService",
			Instances:            domain.Servers{servers[2]},
			NumHealthy:           0,
			AverageCPUPercent:    95.00,
			AverageMemoryPercent: 65.00,
		},
	})
}
