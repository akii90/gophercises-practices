package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
	for i, question := range q.questionsList {
		var userAnswer string
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
	questionsList := make([]question, 0, len(records))
	for index, record := range records {
		question, err := newQuestions(record)
		if err != nil {
			fmt.Printf("Question in line%d is invlaid: %s\n", index+1, err)
			continue
		}
		questionsList = append(questionsList, question)
	}
	return &quiz{
		questionsList: questionsList,
		userScore:     0,
		fullScore:     len(questionsList),
	}
}

func newQuestions(row []string) (question, error) {
	var err error
	text := strings.TrimSpace(row[0])
	answer := strings.TrimSpace(row[1])
	if text == "" || answer == "" {
		err = errors.New("question or answer can not be empty")
	}
	return question{
		text:   text,
		answer: answer,
	}, err
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
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	// Todo， implement shuffle flag, shuffle the quiz order each time it is run
	flag.Parse()

	records := readCsvFile(file)
	quizGame := newQuiz(records)

	ch := make(chan struct{})
	defer close(ch)

	go func() {
		quizGame.challenge()
		ch <- struct{}{}
	}()

	// time limit
	select {
	case <-ch:
		quizGame.scoreOutput()
	case <-time.After(time.Duration(*timeLimit) * time.Second):
		quizGame.scoreOutput()
	}
}
