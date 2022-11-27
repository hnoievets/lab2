package pipeline

type Pipeline struct {
	reverseFlag  string
	NumSortField int
	headerFlag   string
	header       string //заголовок для усіх файлів
}

// функція-конструктор
func NewPipeline(reverseFlag string, NumSortField int, headerFlag string) *Pipeline {
	return &Pipeline{
		reverseFlag:  reverseFlag,
		NumSortField: NumSortField,
		headerFlag:   headerFlag,
		header:       "",
	}
}

// функція для запуску Рipeline з читанням з каталогу
func (p *Pipeline) RunWhithDir(directoryFlag string, outFileName string) {

	fnChan := p.readDir(directoryFlag)
	contentChan := p.filesReadingStage(fnChan)
	p.sortContent(contentChan, outFileName)

}

// функція для запуску Рipeline з читанням з файлу
func (p *Pipeline) RunWhithFile(inputFileName string, outFileName string) {
	fnChan := p.readFileName(inputFileName)
	contentChan := p.filesReadingStage(fnChan)
	p.sortContent(contentChan, outFileName)
}
