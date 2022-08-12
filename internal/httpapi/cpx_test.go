package httpapi_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
	"github.com/jace-ys/thought-machine-interview/internal/httpapi"
)

func TestListServers(t *testing.T) {
	tt := []struct {
		Name             string
		DownstreamStatus int
		ExpectedBody     []net.IP
		ExpectedError    error
	}{
		{
			Name:             "Returns servers on response status 200",
			DownstreamStatus: http.StatusOK,
			ExpectedBody: []net.IP{
				net.ParseIP("10.58.1.1"),
				net.ParseIP("10.58.1.2"),
				net.ParseIP("10.58.1.3"),
				net.ParseIP("10.58.1.4"),
				net.ParseIP("10.58.1.5"),
			},
		},
		{
			Name:             "Returns ErrDownstreamUnavailable on response status 500",
			DownstreamStatus: http.StatusInternalServerError,
			ExpectedError:    httpapi.ErrDownstreamUnavailable,
		},
		{
			Name:             "Returns ErrStatusCodeUnknown on unrecognised response status",
			DownstreamStatus: http.StatusUnauthorized,
			ExpectedError:    httpapi.ErrStatusCodeUnknown,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			handler, client := setupCPX(t)

			fixture, err := os.ReadFile("fixtures/list-servers.json")
			assert.NoError(t, err)

			handler.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.DownstreamStatus)
				if tc.DownstreamStatus == 200 {
					w.Write(fixture)
				}
			})

			servers, err := client.ListServers(context.Background())

			if tc.ExpectedError != nil {
				assert.ErrorIs(t, err, tc.ExpectedError)
				assert.Nil(t, servers)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.ExpectedBody, servers)
			}
		})
	}
}

func TestGetServer(t *testing.T) {
	server := &domain.Server{
		IP:            net.ParseIP("10.58.1.1"),
		Service:       "UserService",
		CPUPercent:    61,
		MemoryPercent: 4,
	}

	tt := []struct {
		Name             string
		DownstreamStatus int
		ExpectedBody     *domain.Server
		ExpectedError    error
	}{
		{
			Name:             "Returns server on response status 200",
			DownstreamStatus: http.StatusOK,
			ExpectedBody:     server,
		},
		{
			Name:             "Returns domain.ErrInvalidServerIP on response status 400",
			DownstreamStatus: http.StatusBadRequest,
			ExpectedError:    domain.ErrInvalidServerIP,
		},
		{
			Name:             "Returns ErrDownstreamUnavailable on response status 500",
			DownstreamStatus: http.StatusInternalServerError,
			ExpectedError:    httpapi.ErrDownstreamUnavailable,
		},
		{
			Name:             "Returns ErrStatusCodeUnknown on unrecognised response status",
			DownstreamStatus: http.StatusUnauthorized,
			ExpectedError:    httpapi.ErrStatusCodeUnknown,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			handler, client := setupCPX(t)

			fixture, err := os.ReadFile("fixtures/get-server.json")
			assert.NoError(t, err)

			handler.HandleFunc("/10.58.1.1", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.DownstreamStatus)
				if tc.DownstreamStatus == 200 {
					w.Write(fixture)
				}
			})

			server, err := client.GetServer(context.Background(), net.ParseIP("10.58.1.1"))

			if tc.ExpectedError != nil {
				assert.ErrorIs(t, err, tc.ExpectedError)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedBody, server)
			}
		})
	}
}

func setupCPX(t *testing.T) (*http.ServeMux, domain.MonitoringService) {
	handler := http.NewServeMux()

	server := httptest.NewServer(handler)
	serverURL, err := url.Parse(server.URL)
	assert.NoError(t, err)

	client := httpapi.NewCPXClient()
	client.BaseURL = serverURL

	return handler, client
}
