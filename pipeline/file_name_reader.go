package pipeline

//функція для передачі імені файлу у канал
func (p *Pipeline) readFileName(fileName string) (fnameCh chan string) {
	fnameCh = make(chan string)

	go func() {
		fnameCh <- fileName
		close(fnameCh)
	}()
	return fnameCh
}
