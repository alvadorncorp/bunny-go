package api

import (
	"context"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/alvadorncorp/bunny-go/internal/logger"
)

type Client interface {
	storage.Client
}

type bunny struct {
	storageClient storage.Client
}

type ClientParams struct {
	StorageName      string
	StorageEndpoint  string
	StorageAccessKey string
	APIKey           string
}

type optionalsParams struct {
	logger logger.Logger
}

type Option = func(o *optionalsParams)

func WithLogger(logger logger.Logger) Option {
	return func(o *optionalsParams) {
		o.logger = logger
	}
}

func New(params ClientParams, opts ...Option) (Client, error) {
	optParams := optionalsParams{}

	for _, apply := range opts {
		apply(&optParams)
	}

	storageClient, err := storage.New(
		storage.ClientParams{
			StorageEndpoint: params.StorageEndpoint,
			StorageName:     params.StorageName,
			APIKey:          params.StorageAccessKey,
		}, storage.WithLogger(optParams.logger))

	if err != nil {
		return nil, err
	}

	return &bunny{
		storageClient: storageClient,
	}, nil
}

func (b bunny) UploadFile(ctx context.Context, file *storage.File) error {
	return b.storageClient.UploadFile(ctx, file)
}
