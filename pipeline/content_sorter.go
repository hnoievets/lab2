package pipeline

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// функція сортування та виводу контенту з каналу
func (p *Pipeline) sortContent(content chan []string, outFileName string) {

	var buff = make([][]string, 0, 1000)

	//додавання до буферу с каналу
	for line := range content {
		buff = append(buff, line)
	}

	//соттування оберненне або звичайне
	if p.reverseFlag == "" {
		sort.Slice(buff, func(i, j int) bool { return buff[i][p.NumSortField-1] > buff[j][p.NumSortField-1] })
	} else {
		sort.Slice(buff, func(i, j int) bool { return buff[i][p.NumSortField-1] < buff[j][p.NumSortField-1] })
	}

	//виведення у файл
	if outFileName != "" {
		fileOut, err := os.OpenFile(outFileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)

		if err != nil {
			log.Fatal(err)
		}

		defer fileOut.Close()

		//виведення заголовку
		if p.header != "" {
			_, err = fileOut.WriteString(strings.Join(strings.Split(p.header, ","), ""))
			fileOut.WriteString("\n")

			if err != nil {
				log.Fatal(err)
			}
		}

		//виведення з буферу
		for _, value := range buff {
			_, err = fileOut.WriteString(strings.Join(value, ""))
			fileOut.WriteString("\n")

			if err != nil {
				log.Fatal(err)
			}
		}

	} else { //виведення у консоль
		fmt.Printf("===Result===\n\n")

		//виведення заголовку
		if p.header != "" {
			fmt.Printf("%s\n", strings.Join(strings.Split(p.header, ","), ""))
		}

		//виведення з буферу
		for i := range buff {
			fmt.Printf("%v\n", strings.Join(buff[i], ""))
		}

	}
}
