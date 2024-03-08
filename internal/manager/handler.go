package manager

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/alvadorncorp/bunny-go/internal/bunny/storage"
)

type fileContainer struct {
	filename string
	filepath string
}

func readDirFiles(dirpath string, ch chan fileContainer) error {
	defer close(ch)
	if err := readDirFiles2(dirpath, "", ch); err != nil {
		log.Println("failure", err)
		return err
	}

	return nil
}

func readDirFiles2(basepath, aggregatedPath string, ch chan fileContainer) error {
	pathToSearch := basepath
	if aggregatedPath != "" {
		pathToSearch = basepath + "/" + aggregatedPath
	}
	entries, err := os.ReadDir(pathToSearch)
	if err != nil {
		log.Println("failure", err)
		return err
	}

	for _, entry := range entries {
		path := aggregatedPath + "/" + entry.Name()
		if aggregatedPath == "" {
			path = entry.Name()
		}
		if entry.IsDir() {
			err := readDirFiles2(basepath, path, ch)
			if err != nil {
				log.Println("failure", err)
				return err
			}
		} else {
			log.Println("not dir", entry.Name())
			ch <- fileContainer{
				filename: entry.Name(),
				filepath: path,
			}
		}
	}

	return nil
}

func (m *manager) Upload(ctx context.Context, args UploadArgs) error {
	log.Println("start reading files")
	ch := make(chan fileContainer, 8)
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
			go func(forFile fileContainer) error {
				defer wg.Done()
				f, err := os.Open(args.SourcePath + "/" + forFile.filepath)
				if err != nil {
					log.Println("file cant be opened", forFile.filepath)
					return err
				}

				defer f.Close()

				if err = m.bunny.UploadFile(
					ctx, &storage.File{
						Buffer:          f,
						Filename:        forFile.filepath,
						DestinationPath: args.DestinationPath,
						ContentType:     "",
						CacheControl:    args.CacheControl,
					}); err != nil {
					log.Println("failure", err)
					return err
				}
				return nil
			}(file)

		}
	}

	wg.Wait()

	return nil
}
