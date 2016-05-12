package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	gconn "github.com/cloudfoundry-incubator/garden/client/connection"
	"github.com/cloudfoundry-incubator/garden/routes"
	"github.com/concourse/atc/db"
	"github.com/concourse/baggageclaim/client"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

//go:generate counterfeiter . GardenConnectionFactoryDB
type GardenConnectionFactoryDB interface {
	GetWorker(string) (db.SavedWorker, bool, error)
}

//go:generate counterfeiter . GardenConnectionFactory
type GardenConnectionFactory interface {
	BuildConnection() gconn.Connection
	BuildConnectionFromDB() (gconn.Connection, error)
	GetDB() GardenConnectionFactoryDB
	CreateRetryableHttpClient() http.Client
}

type gardenConnectionFactory struct {
	db         GardenConnectionFactoryDB
	streamer   WorkerHijackStreamer
	logger     lager.Logger
	workerName string
	address    string
}

func NewGardenConnectionFactory(
	db GardenConnectionFactoryDB,
	logger lager.Logger,
	workerName string,
	address string,
) GardenConnectionFactory {
	return &gardenConnectionFactory{
		db:         db,
		logger:     logger,
		workerName: workerName,
		address:    address,
	}
}

func (gcf *gardenConnectionFactory) BuildConnection() gconn.Connection {
	return gcf.buildConnection(gcf.address)
}

func (gcf *gardenConnectionFactory) BuildConnectionFromDB() (gconn.Connection, error) {
	savedWorker, found, err := gcf.db.GetWorker(gcf.workerName)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrMissingWorker
	}

	return gcf.buildConnection(savedWorker.GardenAddr), nil
}

func (gcf *gardenConnectionFactory) GetDB() GardenConnectionFactoryDB {
	return gcf.db
}

type WorkerHijackStreamer struct {
	HijackStreamer gconn.HijackStreamer
	HttpClient     http.Client
	req            *rata.RequestGenerator
}

func (h *WorkerHijackStreamer) Stream(handler string, body io.Reader, params rata.Params, query url.Values, contentType string) (io.ReadCloser, error) {
	request, err := h.req.CreateRequest(handler, params, body)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	if query != nil {
		request.URL.RawQuery = query.Encode()
	}

	httpResp, err := h.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		defer httpResp.Body.Close()

		var result garden.Error
		err := json.NewDecoder(httpResp.Body).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("bad response: %s", err)
		}

		return nil, result.Err
	}

	return httpResp.Body, nil
}

func (h *WorkerHijackStreamer) Hijack(handler string, body io.Reader, params rata.Params, query url.Values, contentType string) (net.Conn, *bufio.Reader, error) {
	return h.HijackStreamer.Hijack(handler, body, params, query, contentType)
}

// TODO: move this up
type WorkerLookupRoundTripper struct {
	gcf              *gardenConnectionFactory
	httpRoundTripper client.RoundTripper
}

func (roundTripper *WorkerLookupRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	// TODO: alternatively...workerName := request.URL.Host
	savedWorker, found, err := (*roundTripper.gcf).GetDB().GetWorker(roundTripper.gcf.workerName)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrMissingWorker
	}

	(*request.URL).Host = savedWorker.GardenAddr
	return roundTripper.httpRoundTripper.RoundTrip(request)
}

func (gcf *gardenConnectionFactory) CreateRetryableHttpClient() http.Client {
	retryRoundTripper := client.RetryRoundTripper{
		Logger:  lager.NewLogger("retryable-http-client"),
		Sleeper: clock.NewClock(),
		RetryPolicy: ExponentialRetryPolicy{
			Timeout: 60 * time.Minute,
		},
		RoundTripper: &WorkerLookupRoundTripper{
			httpRoundTripper: &http.Transport{
				DisableKeepAlives: true,
			},
			gcf: gcf,
		},
	}

	return http.Client{
		Transport: retryRoundTripper.RoundTripper,
	}
}

func (gcf *gardenConnectionFactory) buildConnection(address string) gconn.Connection {
	hijacker := WorkerHijackStreamer{
		HijackStreamer: gconn.NewHijackStreamer("tcp", address),
		HttpClient:     gcf.CreateRetryableHttpClient(),
		req:            rata.NewRequestGenerator("http://"+address, routes.Routes),
	}
	return gconn.NewWithHijacker(&hijacker, gcf.logger)
}
