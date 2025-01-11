package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// flag for command line
	file := flag.String("file", "problems.csv", "a csv file in a format of 'question,answer'")
	//limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// open file
	quizFile, err := os.Open(*file)
	if err != nil {
		fmt.Printf("Open file error:\n  %v", err)
		return
	}
	defer quizFile.Close()

	// read csv file
	r := csv.NewReader(quizFile)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Read csv error:\n  %v", err)
		return
	}

	// start quiz
	score := 0
	fullScore := len(records)
	var userAnswer string
	for index, record := range records {
		question := record[0]
		correctAnswer := record[1]
		fmt.Printf("Problem #%d: %v = ", index+1, question)
		fmt.Scanln(&userAnswer)
		if strings.EqualFold(userAnswer, correctAnswer) {
			score++
		}
	}
	fmt.Printf("\nYou scored %d of %d.", score, fullScore)
}
