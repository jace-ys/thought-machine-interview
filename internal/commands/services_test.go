package commands

import (
	"bytes"
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
	"github.com/jace-ys/thought-machine-interview/internal/domain/domainfakes"
)

func TestServicesCommand(t *testing.T) {
	server1 := &domain.Server{
		IP:            net.ParseIP("10.58.1.1"),
		Service:       "UserService",
		CPUPercent:    20.00,
		MemoryPercent: 60.00,
	}
	server2 := &domain.Server{
		IP:            net.ParseIP("10.58.1.2"),
		Service:       "StorageService",
		CPUPercent:    100.00,
		MemoryPercent: 35.00,
	}
	server3 := &domain.Server{
		IP:            net.ParseIP("10.58.1.3"),
		Service:       "GeoService",
		CPUPercent:    75.00,
		MemoryPercent: 95.00,
	}

	tt := []struct {
		Name            string
		SetupFake       func(fake *domainfakes.FakeMonitoringService)
		ExpectedOutputs []string
		ExpectedError   bool
	}{
		{
			Name: "Correctly prints services info",
			SetupFake: func(fake *domainfakes.FakeMonitoringService) {
				fake.ListServersReturns([]net.IP{server1.IP, server2.IP, server3.IP}, nil)
				fake.GetServerReturnsOnCall(0, server1, nil)
				fake.GetServerReturnsOnCall(1, server2, nil)
				fake.GetServerReturnsOnCall(2, server3, nil)
			},
			ExpectedOutputs: []string{
				`GeoService\s+\|\s+0/1\s+\|\s+75.00%\s+\|\s+95.00%`,
				`StorageService\s+\|\s+0/1\s+\|\s+100.00%\s+\|\s+35.00%`,
				`UserService\s+\|\s+1/1\s+\|\s+20.00%\s+\|\s+60.00%`,
			},
		},
		{
			Name: "Returns an error if ListServers fails",
			SetupFake: func(fake *domainfakes.FakeMonitoringService) {
				fake.ListServersReturns(nil, errors.New("internal server error"))
			},
			ExpectedError: true,
		},
		{
			Name: "Returns an error if GetServer fails",
			SetupFake: func(fake *domainfakes.FakeMonitoringService) {
				fake.ListServersReturns([]net.IP{server1.IP, server2.IP, server3.IP}, nil)
				fake.GetServerReturnsOnCall(0, server1, nil)
				fake.GetServerReturnsOnCall(1, server2, nil)
				fake.GetServerReturnsOnCall(2, nil, errors.New("internal server error"))
			},
			ExpectedError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			service := new(domainfakes.FakeMonitoringService)
			if tc.SetupFake != nil {
				tc.SetupFake(service)
			}

			cmd := Services(service)

			var buf bytes.Buffer
			err := cmd.execute(context.Background(), &buf)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.Empty(t, buf.Bytes())
			} else {
				assert.NoError(t, err)
				for _, re := range tc.ExpectedOutputs {
					assert.Regexp(t, re, buf.String())
				}
			}
		})
	}
}
