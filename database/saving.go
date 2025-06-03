package database

import (
	"io"
	"os"
)

func UniqueSave(filename string) error {
	sourceFileF, err := os.Open("vault/fibonacciNumbers.txt")
	if err != nil {
		return err
	}
	sourceFilePN, err := os.Open("vault/primeNumbers.txt")
	if err != nil {
		return err
	}

	generatedName := filename
	folderPath := "database/"

	destinationFileF, err := os.Create(folderPath + generatedName + "f" + ".txt")
	if err != nil {
		return err
	}
	destinationFilePN, err := os.Create(folderPath + generatedName + "pn" + ".txt")
	if err != nil {
		return err
	}

	_, err = io.Copy(destinationFileF, sourceFileF)
	if err != nil {
		return err
	}

	_, err = io.Copy(destinationFilePN, sourceFilePN)
	if err != nil {
		return err
	}

	return nil

}
