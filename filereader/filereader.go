package filereader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/widget"
)

func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("ошибка открытия файла: %s", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordCounter := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		wordCounter += len(words)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("ошибка сканирования файла: %s", filename)
	}
	return wordCounter, nil
}

func ReadFile(filename string, output *widget.Label) {
	data, err := os.ReadFile(filename)
	if err != nil {
		output.SetText(fmt.Sprintf("Ошибка чтения файла: %s", filename))
		return
	}

	var whichFile string
	if filename == "primeNumbers.txt" {
		whichFile = "простых чисел."
	} else {
		whichFile = "чисел Фибоначчи."
	}

	output.SetText(fmt.Sprintf("%s\n<-- Список %s", string(data), whichFile))
}
