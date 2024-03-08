package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/alvadorncorp/bunny-go/api"
	"github.com/alvadorncorp/bunny-go/internal/manager"
	"github.com/spf13/cobra"
)

func uploadCmd(storageName, storageEndpoint, storageAPIKey *string) *cobra.Command {
	var source, destination, pattern, cacheControl, contentEncoding string

	upload := &cobra.Command{
		Use: "upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			var re *regexp.Regexp
			if strings.TrimSpace(pattern) != "" {
				compiledRegexp, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				re = compiledRegexp
			}

			bunnyClient, err := api.New(
				api.ClientParams{
					StorageName:     *storageName,
					StorageEndpoint: *storageEndpoint,
					StorageKey:      *storageAPIKey,
					APIKey:          "",
				})

			if err != nil {
				return err
			}

			m := manager.New(bunnyClient)
			return m.Upload(
				cmd.Context(), manager.UploadArgs{
					Pattern:         re,
					SourcePath:      source,
					DestinationPath: destination,
					CacheControl:    cacheControl,
				})

		},
	}

	flags := upload.Flags()
	flags.StringVarP(&source, "source-path", "", "", "source path")
	flags.StringVarP(&destination, "destination-path", "", "", "destination path")
	flags.StringVarP(&pattern, "pattern", "", "", "regex pattern for files")
	flags.StringVarP(&cacheControl, "cache-control", "", "max-age=2592000", "cache control header")
	flags.StringVarP(&contentEncoding, "content-encoding", "", "", "content encoding header")

	return upload
}

func main() {
	var storageName, storageEndpoint, storageApiKey string
	cmd := &cobra.Command{
		Use: "api",
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&storageName, "storage-name", "", "", "storage name")
	flags.StringVarP(&storageEndpoint, "storage-endpoint", "", "br", "storage endpoint")
	flags.StringVarP(&storageApiKey, "storage-api-key", "", "", "storage api key")

	cmd.AddCommand(uploadCmd(&storageName, &storageEndpoint, &storageApiKey))

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
