package pipeline

import (
	"log"
	"os"
	"path/filepath"
)

// функція для читання каталогу з обходом, передає імена файлів у канал
func (p *Pipeline) readDir(dir string) (fnames chan string) {
	fnames = make(chan string)

	go func(dir string) {
		defer close(fnames)

		listFileNames := []string{}
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				listFileNames = append(listFileNames, path)
				return nil
			})
		if err != nil {
			log.Println(err)
		}

		for _, f := range listFileNames {
			fnames <- f
		}
	}(dir)

	return fnames
}
