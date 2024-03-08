package manager

import (
	"regexp"

	"github.com/alvadorncorp/bunny-go/api"
	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
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
