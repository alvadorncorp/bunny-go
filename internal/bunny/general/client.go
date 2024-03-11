package general

import (
	"context"
	"net/http"

	"github.com/alvadorncorp/bunny-go/internal/bunny"
	"github.com/alvadorncorp/bunny-go/internal/logger"
)

const (
	defaultBufferSize = 4 * 1024 // 4 KiB
	bunnyAPIUrl       = "https://api.bunny.net"
)

type Option = func(b *bunnyClient)

type Client interface {
	ListPullZones(ctx context.Context, params ListPullZoneParams) (bunny.Page[PullZone], error)
}

type BunnyClientParams struct {
	APIKey string
}

type bunnyClient struct {
	baseAPIUrl string
	client     *http.Client
	apiKey     string
	logger     logger.Logger
}

func WithHttpClient(hc *http.Client) Option {
	return func(b *bunnyClient) {
		b.client = hc
	}
}

func WithLogger(log logger.Logger) Option {
	return func(b *bunnyClient) {
		b.logger = log
	}
}

func New(params BunnyClientParams, options ...Option) Client {
	client := &bunnyClient{
		baseAPIUrl: bunnyAPIUrl,
		client:     &http.Client{},
		apiKey:     params.APIKey,
		logger:     logger.NewMockLogger(),
	}

	for _, applyOption := range options {
		applyOption(client)
	}

	return client
}
