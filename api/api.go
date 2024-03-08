package api

import (
	"context"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
)

type Client interface {
	storage.Client
}

type bunny struct {
	storageClient storage.Client
}

type ClientParams struct {
	StorageName     string
	StorageEndpoint string
	StorageKey      string
	APIKey          string
}

func New(params ClientParams) (Client, error) {
	storageClient, err := storage.New(
		storage.ClientParams{
			StorageEndpoint: params.StorageEndpoint,
			StorageName:     params.StorageName,
			APIKey:          params.APIKey,
		})

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
