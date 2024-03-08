package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

type Option = func(sc *storageClient)

type StorageClientParams struct {
	StorageEndpoint string
	StorageName     string
	APIKey          string
}

func WithHttpClient(hc *http.Client) Option {
	return func(sc *storageClient) {
		sc.client = hc
	}
}

func New(params StorageClientParams, options ...Option) Client {
	sc := &storageClient{
		client:     &http.Client{},
		apiKey:     params.APIKey,
		baseAPIUrl: fmt.Sprint("https://", params.StorageEndpoint, storageApiUrl, "/", params.StorageName, "/"),
	}

	for _, applyOption := range options {
		applyOption(sc)
	}

	return sc
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, f.Buffer)
	if err != nil {
		log.Println("failure", path, err)
		return err
	}
	log.Println("uploading ", destination)

	headers := req.Header
	headers.Add("content-type", f.ContentType)
	headers.Add("AccessKey", b.apiKey)
	headers.Add("cache-control", f.CacheControl)

	res, err := b.client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	if res.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("failure", body)
		return err
	}

	return fmt.Errorf(string(body))
}
