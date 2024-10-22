package storage_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorageClient_UploadFile(t *testing.T) {
	emptyCtx := context.Background()

	t.Run("file has been uploaded successfully", func(t *testing.T) {
		file := &storage.LocalFile{
			Buffer:          bytes.NewBuffer([]byte{}),
			DestinationPath: "path",
			Filename:        "filename.jpg",
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/path/filename.jpg", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		srv := httptest.NewServer(mux)

		strClient := storage.New(storage.ClientParams{
			StorageName: "name",
			APIKey:      "api-key",
		}, storage.WithTestUrl(srv.URL))

		err := strClient.UploadFile(emptyCtx, file)
		defer srv.Close()
		assert.NoError(t, err)
	})
}
