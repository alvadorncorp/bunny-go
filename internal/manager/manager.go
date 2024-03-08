package manager

import (
	"regexp"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/alvadorncorp/bunny-go/pkg/api"
)

type UploadArgs struct {
	Pattern         *regexp.Regexp
	SourcePath      string
	DestinationPath string
	CacheControl    string
	ContentEncoding string
}

type manager struct {
	bunny storage.Client
}

func New(bunny api.Client) *manager {
	return &manager{bunny: bunny}
}
