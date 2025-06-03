package primenumbers

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

func PrimeNumbersFounder(limiter int, filename string, wg *sync.WaitGroup, resultsChan chan<- string) {
	defer wg.Done()

	file, err := os.Create(filename)
	if err != nil {
		resultsChan <- "Ошибка при создании файла primeNumbers.txt"
		return
	}
	defer file.Close()

	counter := 1
	num := 2
	for counter < limiter {
		if isPrime(num) {
			file.WriteString(strconv.Itoa(num) + " ")
			resultsChan <- fmt.Sprintf("[Поиск простых чисел] | Шаг: %d | Найденное простое число: %d", counter, num)
			counter++
		}
		num++
	}
}

func isPrime(numToAnalyze int) bool {
	var i float64
	if numToAnalyze < 2 {
		return false
	}
	for i = 2; i <= math.Sqrt(float64(numToAnalyze)); i++ {
		if numToAnalyze%int(i) == 0 {
			return false
		}
	}
	return true
}
