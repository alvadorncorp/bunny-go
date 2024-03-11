package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alvadorncorp/bunny-go/internal/logger"
)

const (
	binaryContentType = "application/octet-stream"
	storageApiUrl     = "storage.bunnycdn.com"
)

type Client interface {
	UploadFile(ctx context.Context, file *File) error
}

type storageClient struct {
	baseAPIUrl string
	client     *http.Client
	apiKey     string
	logger     logger.Logger
}

type Option = func(sc *storageClient)

type ClientParams struct {
	StorageEndpoint string
	StorageName     string
	APIKey          string
}

func WithHttpClient(hc *http.Client) Option {
	return func(sc *storageClient) {
		sc.client = hc
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(sc *storageClient) {
		sc.logger = logger
	}
}

func New(params ClientParams, options ...Option) (Client, error) {
	url := fmt.Sprint("https://", params.StorageEndpoint, ".", storageApiUrl, "/", params.StorageName)
	if params.StorageEndpoint == "" {
		url = fmt.Sprint("https://", storageApiUrl, "/", params.StorageName)
	}

	sc := &storageClient{
		client:     &http.Client{},
		apiKey:     params.APIKey,
		baseAPIUrl: url,
		logger:     logger.NewMockLogger(),
	}

	for _, applyOption := range options {
		applyOption(sc)
	}

	sc.logger = sc.logger.With(logger.String("client", "bunny-storage"))

	return sc, nil
}

type File struct {
	Buffer          io.Reader
	DestinationPath string
	Filename        string
	ContentType     string
	CacheControl    string
}

func sanitizeDestinationPath(destinationPath string) string {
	if strings.HasPrefix(destinationPath, "./") {
		return destinationPath[2:]
	}

	if strings.HasPrefix(destinationPath, ".") {
		return destinationPath[1:]
	}

	return destinationPath
}

func (b *storageClient) UploadFile(ctx context.Context, f *File) error {
	destination := f.Filename
	if destinationPath := sanitizeDestinationPath(f.DestinationPath); destinationPath != "" {
		destination = fmt.Sprintf("%s/%s", destinationPath, f.Filename)
	}

	apiPath := fmt.Sprintf("%s/%s", b.baseAPIUrl, destination)
	loggerChild := b.logger.With(
		logger.String("context", "storageClient.uploadFile"),
		logger.String("filename", f.Filename),
		logger.String("api-path", apiPath),
		logger.String("destination-path", destination))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiPath, f.Buffer)
	if err != nil {
		loggerChild.Error(err, "new request creation failure")
		return err
	}

	headers := req.Header
	headers.Set("content-type", binaryContentType)
	headers.Set("AccessKey", b.apiKey)
	headers.Set("cache-control", f.CacheControl)

	loggerChild.Debug("starting upload request...")
	res, err := b.client.Do(req)
	if err != nil {
		loggerChild.Error(err, "request failure")
		return err
	}

	if res.StatusCode < 400 {
		loggerChild.Debug("request completed")
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerChild.Error(err, "failure while reading failure response")
		return err
	}

	return fmt.Errorf(string(body))
}
