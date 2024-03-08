package bunny

import (
	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
)

type Client interface {
	storage.Client
}

type bunny struct {
	storage.Client
}

type ClientParams struct {
	StorageName     string
	StorageEndpoint string
	StorageKey      string
	APIKey          string
}

func New(params ClientParams) Client {
	storageClient := storage.New(
		storage.StorageClientParams{
			StorageEndpoint: params.StorageEndpoint,
			StorageName:     params.StorageName,
			APIKey:          params.APIKey,
		})

	return &bunny{
		storageClient,
	}
}
