package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type quiz struct {
	questionsList []question
	userScore     int
	fullScore     int
}

type question struct {
	text   string
	answer string
}

// challenge give user questions, wait user to answer, recording user score.
func (q *quiz) challenge() {
	var userAnswer string
	for i, question := range q.questionsList {
		number := i + 1
		fmt.Printf("Problem #%d: %v = ", number, question.text)
		fmt.Scanln(&userAnswer)
		if strings.EqualFold(userAnswer, question.answer) {
			q.userScore++
		}
	}
}

// scoreOutput print the score.
func (q *quiz) scoreOutput() {
	fmt.Printf("\nYou scored %d of %d.", q.userScore, q.fullScore)
}

func newQuiz(records [][]string) *quiz {
	questionsList := make([]question, len(records))
	for index, record := range records {
		if len(record) < 2 {
			fmt.Println("Error format for question and answer")
			continue
		}
		questionsList[index] = newQuestions(record)
	}
	return &quiz{
		questionsList: questionsList,
		userScore:     0,
		fullScore:     len(records),
	}
}

func newQuestions(row []string) question {
	return question{
		text:   row[0],
		answer: row[1],
	}
}

// readCsvFile read a csv file, get the record in csv file.
func readCsvFile(file *string) (records [][]string) {
	// open file
	quizFile, err := os.Open(*file)
	if err != nil {
		log.Fatalf("Open file error:\n  %v", err)
	}
	defer quizFile.Close()

	// read csv file
	r := csv.NewReader(quizFile)
	records, err = r.ReadAll()
	if err != nil {
		log.Fatalf("Read csv file error:\n  %v", err)
	}

	return records
}

func main() {
	// flag for command line
	file := flag.String("file", "problems.csv", "a csv file in a format of 'question,answer'")
	// Todo， implement limit flag
	//timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	// Todo， implement shuffle flag, shuffle the quiz order each time it is run
	flag.Parse()

	records := readCsvFile(file)
	quizGame := newQuiz(records)
	quizGame.challenge()
	quizGame.scoreOutput()
}
