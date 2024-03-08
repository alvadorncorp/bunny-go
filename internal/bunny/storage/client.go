package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/alvadorncorp/bunny-go/internal/logger"
)

const (
	storageApiUrl = "storage.bunny.net"
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
	logClient, err := logger.NewZapLogger()
	if err != nil {
		return nil, err
	}

	sc := &storageClient{
		client:     &http.Client{},
		apiKey:     params.APIKey,
		baseAPIUrl: fmt.Sprint("https://", params.StorageEndpoint, storageApiUrl, "/", params.StorageName, "/"),
		logger:     logClient,
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

func (b *storageClient) UploadFile(ctx context.Context, f *File) error {
	destination := f.Filename
	if f.DestinationPath != "" {
		destination = fmt.Sprintf("%s/%s", f.DestinationPath, f.Filename)
	}

	path := fmt.Sprintf("%s/%s", b.baseAPIUrl, destination)

	loggerChild := b.logger.With(logger.String("path", path))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, f.Buffer)
	if err != nil {
		loggerChild.Error(err, "new request failure")
		return err
	}
	loggerChild.Debug("uploading file => start")

	headers := req.Header
	headers.Add("content-type", f.ContentType)
	headers.Add("AccessKey", b.apiKey)
	headers.Add("cache-control", f.CacheControl)

	res, err := b.client.Do(req)
	if err != nil {
		loggerChild.Error(err, "uploading file => failure")
		return err
	}

	if res.StatusCode < 400 {
		loggerChild.Debug("uploading file => completed")
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		loggerChild.Error(err, "failure while reading failure response")
		return err
	}

	return fmt.Errorf(string(body))
}
