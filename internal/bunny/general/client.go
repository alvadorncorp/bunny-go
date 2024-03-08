package general

import (
	"context"
	"net/http"

	"github.com/alvadorncorp/bunny-go/internal/bunny"
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
}

func WithHttpClient(hc *http.Client) Option {
	return func(b *bunnyClient) {
		b.client = hc
	}
}

func New(params BunnyClientParams, options ...Option) Client {
	client := &bunnyClient{
		baseAPIUrl: bunnyAPIUrl,
		client:     &http.Client{},
		apiKey:     params.APIKey,
	}

	for _, applyOption := range options {
		applyOption(client)
	}

	return client
}
