package calculating

import (
	"evm/fibonacci"
	"evm/primenumbers"
	"sync"

	"fyne.io/fyne/v2/widget"
)

func FAndPfoundind(fiblimiter int, primelimiter int, resultsChan chan<- string, startButton *widget.Button, analyzeButton *widget.Button) {
	var wg sync.WaitGroup
	fibonacciFileName := "vault/fibonacciNumbers.txt"
	primeNumbersFileName := "vault/primeNumbers.txt"
	wg.Add(2)

	go fibonacci.FibonacciFounder(fiblimiter+1, fibonacciFileName, &wg, resultsChan)
	go primenumbers.PrimeNumbersFounder(primelimiter+1, primeNumbersFileName, &wg, resultsChan)

	go func() {
		wg.Wait()          
		close(resultsChan) 
		defer startButton.Enable()
		defer startButton.SetText("Начать поиск")
		defer analyzeButton.Enable()
		
	}()

}
