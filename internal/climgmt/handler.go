package climgmt

import (
	"context"
	"os"
	"sync"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
	"github.com/alvadorncorp/bunny-go/internal/logger"
)

type fileContainer struct {
	filename string
	filepath string
}

func readDirFiles(dirPath string, ch chan fileContainer) error {
	defer close(ch)
	if err := readDirFilesAux(dirPath, "", ch); err != nil {
		return err
	}

	return nil
}

func readDirFilesAux(basePath, aggregatedPath string, ch chan fileContainer) error {
	pathToSearch := basePath
	if aggregatedPath != "" {
		pathToSearch = basePath + "/" + aggregatedPath
	}
	entries, err := os.ReadDir(pathToSearch)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := aggregatedPath + "/" + entry.Name()
		if aggregatedPath == "" {
			path = entry.Name()
		}
		if entry.IsDir() {
			err := readDirFilesAux(basePath, path, ch)
			if err != nil {
				return err
			}
		} else {
			ch <- fileContainer{
				filename: entry.Name(),
				filepath: path,
			}
		}
	}

	return nil
}

func (m *cliManager) Upload(ctx context.Context, args UploadArgs) error {
	ch := make(chan fileContainer, maxConcurrency)
	go readDirFiles(args.SourcePath, ch)

	wg := &sync.WaitGroup{}
	for {
		if ch == nil {
			break
		}

		select {
		case file, ok := <-ch:
			if !ok {
				ch = nil
				break
			}

			wg.Add(1)
			go func(file fileContainer) error {
				defer wg.Done()
				filepath := args.SourcePath + "/" + file.filepath
				f, err := os.Open(filepath)
				if err != nil {
					m.logger.Error(err, "file can't be open", logger.String("filepath", file.filepath))
					return err
				}

				defer f.Close()

				if err = m.bunny.UploadFile(
					ctx, &storage.LocalFile{
						Buffer:          f,
						Filename:        file.filepath,
						DestinationPath: args.DestinationPath,
					}); err != nil {
					m.logger.Error(err, "upload file failure", logger.String("filename", file.filepath))
					return err
				}
				return nil
			}(file)

		}
	}

	wg.Wait()
	return nil
}
