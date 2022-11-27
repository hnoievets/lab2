package main

import (
	"bufio"
	"flag"
	"fmt"
	"lab2/pipeline"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {

	fmt.Println("===Start===")

	var (
		iFlag = flag.String("i", "", "Use a file with the name file-name as an input")
		oFlag = flag.String("o", "", "Use a file with the name file-name as an output")
		rFlag = flag.String("r", "0", "Sort input lines in reverse order")
		fFlag = flag.Int("f", 1, "Sort input lines by value number N")
		hFlag = flag.String("h", "0", "The first line is a header that must be ignored during sorting but included in the output")
		dFlag = flag.String("d", "", "Use a directory with the subdirectories and files as an input")
	)

	flag.Parse()

	if *iFlag != "" && *dFlag != "" {
		log.Fatalln("ERROR: Using incompatible flags (-i and -d) !!!")
	}

	//створення об'єкту класу Рipeline, якщо читатння з каталогу або файл
	if *iFlag != "" || *dFlag != "" {
		p := pipeline.NewPipeline(*rFlag, *fFlag, *hFlag)

		//запуск Рipeline з читанням з каталогу
		if *dFlag != "" {
			p.RunWhithDir(*dFlag, *oFlag)

			//запуск Рipeline з читанням з файлу
		} else if *iFlag != "" {
			p.RunWhithFile(*iFlag, *oFlag)
		}

	} else {
		// робота з консоллю
		s := bufio.NewScanner(os.Stdin)
		var n int //змінна для перевірки к-сті слів у рядку
		var header string
		content := [][]string{}

		//читання з консолі
		for s.Scan() {
			line := s.Text()

			if line == "" {
				break
			}

			row := strings.Split(line, ",")

			if n == 0 {
				n = len(row)

				if *fFlag > n || *fFlag < 1 {
					log.Fatalf("ERROR: The line has %d fields, but for sorting needs %d\n", n, *fFlag)
				}
			}
			if len(row) != n {
				log.Fatalf("ERROR: The length of line must be %d, but it is %d\n", n, len(row))
			}
			if *hFlag == "" && header == "" {
				header = line
				continue
			}

			content = append(content, row)
		}

		//соттування оберненне або звичайне
		if *rFlag == "" {
			sort.Slice(content, func(i, j int) bool { return content[i][*fFlag-1] > content[j][*fFlag-1] })
		} else {
			sort.Slice(content, func(i, j int) bool { return content[i][*fFlag-1] < content[j][*fFlag-1] })
		}

		//виведення у файл
		if *oFlag != "" {
			//відкриття файлу для читатння та запису, видалення змісту при відкритті, створення
			fileOut, err := os.OpenFile(*oFlag, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)

			if err != nil {
				log.Fatal(err)
			}

			if header != "" {
				_, err = fileOut.WriteString(strings.Join(strings.Split(header, ","), ""))
				fileOut.WriteString("\n")

				if err != nil {
					log.Fatal(err)
				}
			}

			for _, value := range content {
				_, err = fileOut.WriteString(strings.Join(value, ""))
				fileOut.WriteString("\n")

				if err != nil {
					log.Fatal(err)
				}
			}

			defer fileOut.Close()
		}

		//вивід у консоль
		fmt.Printf("===Result===\n\n")

		if header != "" {
			fmt.Printf("%s\n", strings.Join(strings.Split(header, ","), ""))
		}

		for _, value := range content {
			fmt.Printf("%s\n", strings.Join(value, ""))
		}
		fmt.Printf("\n===Finish===\n\n")
	}

}
