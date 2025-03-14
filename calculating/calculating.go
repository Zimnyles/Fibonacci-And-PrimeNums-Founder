package calculating

import (
	"evm/fibonacci"
	"evm/primenumbers"
	"sync"
)

func FAndPfoundind(fiblimiter int, primelimiter int, resultsChan chan<- string) {
	var wg sync.WaitGroup
	fibonacciFileName := "vault/fibonacciNumbers.txt"
	primeNumbersFileName := "vault/primeNumbers.txt"
	wg.Add(2)

	// Запускаем горутины для поиска чисел
	go fibonacci.FibonacciFounder(fiblimiter+1, fibonacciFileName, &wg, resultsChan)
	go primenumbers.PrimeNumbersFounder(primelimiter+1, primeNumbersFileName, &wg, resultsChan)

	// Запускаем горутину, которая закроет канал после завершения всех рабочих горутин
	go func() {
		wg.Wait() // Ждем завершения всех горутин
		// close(resultsChan) // Закрываем канал
	}()
}
