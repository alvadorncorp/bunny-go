package main

import (
	"fmt"

	"github.com/alvadorncorp/bunny-go/internal/climgmt"
	"github.com/alvadorncorp/bunny-go/internal/logger"
	"github.com/alvadorncorp/bunny-go/pkg/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

type storageCmdFlags struct {
	Name     string
	Endpoint string
	APIKey   string
}

type uploadCmdFlags struct {
	source          string
	destination     string
	pattern         string
	cacheControl    string
	contentEncoding string
}

func requiredFlagError(flagName string) error {
	return fmt.Errorf("%s is a required flag", flagName)
}

func validateStorageCmdFlags(opts storageCmdFlags) error {
	if opts.Name == "" {
		return requiredFlagError("storage-name")
	}

	if opts.APIKey == "" {
		return requiredFlagError("storage-access-key")
	}

	return nil
}

func validateUploadCmdOptions(opts uploadCmdFlags) error {
	if opts.source == "" {
		return requiredFlagError("source-path")
	}

	if opts.destination == "" {
		return requiredFlagError("destination-path")
	}

	return nil
}

func uploadCmd(storageFlags *storageCmdFlags) *cobra.Command {
	var options uploadCmdFlags

	upload := &cobra.Command{
		Use:   "upload",
		Short: "upload files",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateUploadCmdOptions(options); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logger.NewZapLogger(zapcore.InfoLevel)

			if err != nil {
				return err
			}

			bunnyClient, err := api.New(
				api.ClientParams{
					StorageName:      storageFlags.Name,
					StorageEndpoint:  storageFlags.Endpoint,
					StorageAccessKey: storageFlags.APIKey,
					APIKey:           "",
				}, api.WithLogger(log))

			if err != nil {
				return err
			}

			m := climgmt.New(bunnyClient, log)

			return m.Upload(
				cmd.Context(), climgmt.UploadArgs{
					SourcePath:      options.source,
					DestinationPath: options.destination,
				})
		},
	}

	flags := upload.Flags()
	flags.StringVarP(&options.source, "source-path", "", "", "source path")
	flags.StringVarP(&options.destination, "destination-path", "", "", "destination path")

	return upload
}

func storageCmd() *cobra.Command {
	var flags storageCmdFlags

	cmd := &cobra.Command{
		Use:   "storage",
		Short: "storage <subcommand>",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateStorageCmdFlags(flags); err != nil {
				return err
			}

			return nil
		},
	}

	pFlags := cmd.PersistentFlags()
	pFlags.StringVarP(&flags.Name, "storage-name", "", "", "storage name")
	pFlags.StringVarP(&flags.Endpoint, "storage-endpoint", "", "", "storage endpoint (e.g.: br,de,us)")
	pFlags.StringVarP(&flags.APIKey, "storage-access-key", "", "", "storage api key")

	cmd.AddCommand(uploadCmd(&flags))
	return cmd
}
