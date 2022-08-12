package httpapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
)

var _ domain.MonitoringService = (*CPXClient)(nil)

type CPXClient struct {
	BaseURL *url.URL
	client  *http.Client
}

func NewCPXClient() *CPXClient {
	u, err := url.Parse("http://localhost:8000/")
	if err != nil {
		panic(err)
	}

	return &CPXClient{
		BaseURL: u,
		client:  http.DefaultClient,
	}
}

type ListServersResponse []net.IP

func (c *CPXClient) ListServers(ctx context.Context) ([]net.IP, error) {
	endpoint := "/servers"
	req, err := NewRequest(ctx, c.BaseURL, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var res ListServersResponse
	rsp, err := Do(c.client, req, &res)
	if err != nil {
		return nil, err
	}

	switch {
	case rsp.StatusCode == http.StatusOK:
		// OK
	case 500 <= rsp.StatusCode && rsp.StatusCode <= 599:
		return nil, fmt.Errorf("%w: %s", ErrDownstreamUnavailable, rsp.HTTPErrorBody)
	default:
		return nil, fmt.Errorf("%w: %d", ErrStatusCodeUnknown, rsp.StatusCode)
	}

	return res, nil
}

func (c *CPXClient) GetServer(ctx context.Context, serverIP net.IP) (*domain.Server, error) {
	endpoint := fmt.Sprintf("/%s", serverIP.String())
	req, err := NewRequest(ctx, c.BaseURL, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var res server
	rsp, err := Do(c.client, req, &res)
	if err != nil {
		return nil, err
	}

	switch {
	case rsp.StatusCode == http.StatusOK:
		// OK
	case rsp.StatusCode == http.StatusBadRequest:
		return nil, domain.ErrInvalidServerIP
	case 500 <= rsp.StatusCode && rsp.StatusCode <= 599:
		return nil, fmt.Errorf("%w: %s", ErrDownstreamUnavailable, rsp.HTTPErrorBody)
	default:
		return nil, fmt.Errorf("%w: %d", ErrStatusCodeUnknown, rsp.StatusCode)
	}

	return res.toDomain(serverIP)
}

type server struct {
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	Service string `json:"service"`
}

func (s *server) toDomain(serverIP net.IP) (*domain.Server, error) {
	cpu, err := parsePercentValue(s.CPU)
	if err != nil {
		return nil, fmt.Errorf("cpu: %s", err)
	}

	memory, err := parsePercentValue(s.Memory)
	if err != nil {
		return nil, fmt.Errorf("memory: %s", err)
	}

	return &domain.Server{
		IP:            serverIP,
		Service:       s.Service,
		CPUPercent:    cpu,
		MemoryPercent: memory,
	}, nil
}

func parsePercentValue(s string) (float64, error) {
	val, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing percent value: %w", err)
	}
	return val, nil
}
