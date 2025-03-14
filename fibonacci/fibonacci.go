package fibonacci

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func FibonacciFounder(limiter int, filename string, wg *sync.WaitGroup, resultsChan chan<- string) {
	defer wg.Done()

	file, err := os.Create(filename)
	if err != nil {
		resultsChan <- "Ошибка при создании файла fibonacciNumbers.txt"
		return
	}

	defer file.Close()

	a, b := 0, 1
	for i := 1; i < limiter; i++ {
		file.WriteString(strconv.Itoa(a) + " ")
		resultsChan <- fmt.Sprintf("[Поиск чисел Фибоначчи] | Шаг: %d | Сумма чисел: %d", i, a+b)
		a, b = b, a+b
	}
}
