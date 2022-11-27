package pipeline

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

// функція для читання контенту фалів у канал, з каналу який містить ім'я файлів
func (p *Pipeline) readFiles(fnames chan string, lines chan []string) {

	var n int         // змінан для перевірки к-сті слів у рядку
	var header string // змінна для заголоку кожного файлу

	go func() {
		defer close(lines)

		for fname := range fnames {

			fileIn, err := os.Open(fname)

			if err != nil {
				log.Fatal(err)
			}

			s := bufio.NewScanner(fileIn)

			//читання файлу
			for s.Scan() {
				line := s.Text()

				if line == "" {
					break
				}

				row := strings.Split(line, ",")

				if n == 0 {
					n = len(row)
					//перевірка корекності сортування за полем
					if p.NumSortField > n || p.NumSortField < 1 {
						log.Fatalf("ERROR: The line has %d fields, but for sorting needs %d\n", n, p.NumSortField)
					}
				}

				//помилка якщо кількість слів у рядках різна
				if len(row) != n {
					log.Fatalf("ERROR: The length of line must be %d, but it is %d - %v\n", n, len(row), fname)
				}

				//читання заголвку файлу без його додавання до контенту
				if p.headerFlag == "" && header == "" {
					header = line
					if p.header == "" {
						p.header = line
					}

					continue

					//помилка якщо заголовки різні
				} else if header != "" && p.header != "" && header != p.header {
					log.Fatalf("ERROR: file %s has a differnt header\n", fname)
				}

				lines <- row
			}

			header = ""
			fileIn.Close()
		}

	}()

}

// функція для конкурентного отримання контенут з каналу
func (p *Pipeline) filesReadingStage(fnames chan string) (allLines chan []string) {

	n := 3 //кількість каналів для отримання контенту

	allLines = make(chan []string)

	lines := make([]chan []string, n)

	//ствоерння 3х каналів для читання
	for i := 0; i < n; i++ {
		lines[i] = make(chan []string)
		p.readFiles(fnames, lines[i])
	}

	//отримання контенту з очікуванням закриття усіх каналів
	go func() {
		defer close(allLines)
		wg := &sync.WaitGroup{}

		for i := range lines {
			wg.Add(1)
			go func(ch chan []string) {
				defer wg.Done()
				for line := range ch {
					allLines <- line
				}
			}(lines[i])
		}

		wg.Wait()
	}()

	return allLines
}
