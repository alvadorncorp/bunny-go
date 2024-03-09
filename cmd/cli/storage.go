package main

import (
	"regexp"
	"strings"

	"github.com/alvadorncorp/bunny-go/internal/climgmt"
	"github.com/alvadorncorp/bunny-go/pkg/api"
	"github.com/spf13/cobra"
)

type storageFlags struct {
	Name     string
	Endpoint string
	APIKey   string
}

func uploadCmd(storageFlags storageFlags) *cobra.Command {
	var source, destination, pattern, cacheControl, contentEncoding string

	upload := &cobra.Command{
		Use:   "upload",
		Short: "upload files",
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
					StorageName:     storageFlags.Name,
					StorageEndpoint: storageFlags.Endpoint,
					StorageKey:      storageFlags.APIKey,
					APIKey:          "",
				})

			if err != nil {
				return err
			}

			m := climgmt.New(bunnyClient)
			return m.Upload(
				cmd.Context(), climgmt.UploadArgs{
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

func storageCmd() *cobra.Command {
	var flags storageFlags

	cmd := &cobra.Command{
		Use:   "storage",
		Short: "storage <subcommand>",
	}

	pFlags := cmd.PersistentFlags()
	pFlags.StringVarP(&flags.Name, "name", "", "", "storage name")
	pFlags.StringVarP(&flags.Endpoint, "endpoint", "", "br", "storage endpoint")
	pFlags.StringVarP(&flags.APIKey, "api-key", "", "", "storage api key")

	cmd.AddCommand(uploadCmd(flags))
	return cmd
}
