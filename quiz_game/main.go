package main

import (
	"fmt"
	"os"
)

func main() {
	// flag for cmd
	quizFile := "problems.csv"

	// open file
	file, err := os.Open(quizFile)
	if err != nil {
		fmt.Printf("Open file error:\n  %v", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()
}
