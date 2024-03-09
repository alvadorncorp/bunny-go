package climgmt

import (
	"context"
	"regexp"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/alvadorncorp/bunny-go/internal/logger"
	"github.com/alvadorncorp/bunny-go/pkg/api"
)

type Manager interface {
	Upload(ctx context.Context, args UploadArgs) error
}

type UploadArgs struct {
	Pattern         *regexp.Regexp
	SourcePath      string
	DestinationPath string
	CacheControl    string
	ContentEncoding string
}

type cliManager struct {
	bunny  storage.Client
	logger logger.Logger
}

func New(bunny api.Client) Manager {
	return &cliManager{bunny: bunny}
}
