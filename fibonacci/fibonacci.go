package fibonacci

import (
	"fmt"
	"os"
	"math/big"

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

    a := big.NewInt(0)
    b := big.NewInt(1)
    sum := new(big.Int)

    for i := 1; i < limiter; i++ {
        file.WriteString(a.String() + " ")
        sum.Add(a, b)
        resultsChan <- fmt.Sprintf("[Поиск чисел Фибоначчи] | Шаг: %d | Сумма чисел: %s", i, sum.String())
        
        a, b = b, sum
    }
}
