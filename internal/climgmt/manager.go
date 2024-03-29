package climgmt

import (
	"context"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/alvadorncorp/bunny-go/internal/logger"
	"github.com/alvadorncorp/bunny-go/pkg/api"
)

const maxConcurrency = 8

type Manager interface {
	Upload(ctx context.Context, args UploadArgs) error
}

type UploadArgs struct {
	SourcePath      string
	DestinationPath string
}

type cliManager struct {
	bunny  storage.Client
	logger logger.Logger
}

func New(bunny api.Client, logger logger.Logger) Manager {
	return &cliManager{bunny: bunny, logger: logger}
}
