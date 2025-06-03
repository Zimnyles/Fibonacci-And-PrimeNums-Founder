package dataanalyze

import (
	"evm/filereader"
	"fmt"
	"os"

	"fyne.io/fyne/v2/widget"
)

func DataAnalyze(output *widget.Label) {
	// Чтение и вывод количества чисел Фибоначчи
	fibonacciNumsCount, err := filereader.CountLines("vault/fibonacciNumbers.txt")
	if err != nil {
		output.SetText("Ошибка чтения файла fibonacciNumbers.txt!")
		return
	}

	// Чтение и вывод количества простых чисел
	primeNumsCount, err := filereader.CountLines("vault/primeNumbers.txt")
	if err != nil {
		output.SetText("Ошибка чтения файла primeNumbers.txt!")
		return
	}

	// Вывод количества чисел
	output.SetText(fmt.Sprintf("Количество чисел Фибоначчи: %d\nКоличество простых чисел: %d", fibonacciNumsCount, primeNumsCount))

	// Чтение и вывод содержимого файлов
	fibonacciData, err := os.ReadFile("vault/fibonacciNumbers.txt")
	if err != nil {
		output.SetText(output.Text + "\nОшибка чтения файла fibonacciNumbers.txt!")
		return
	}

	primeData, err := os.ReadFile("vault/primeNumbers.txt")
	if err != nil {
		output.SetText(output.Text + "\nОшибка чтения файла primeNumbers.txt!")
		return
	}

	// Объединение данных из обоих файлов
	combinedData := fmt.Sprintf(
		"Числа Фибоначчи:\n%s\n\nПростые числа:\n%s",
		string(fibonacciData),
		string(primeData),
	)

	// Вывод объединенных данных
	output.SetText(output.Text + "\n" + combinedData)
}
