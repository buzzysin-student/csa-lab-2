package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// question struct stores a single question and its corresponding answer.
type question struct {
	q, a string
}

type score int

// check handles a potential error.
// It stops execution of the program ("panics") if an error has happened.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// questions reads in questions and corresponding answers from a CSV file into a slice of question structs.
func questions() []question {
	f, err := os.Open("quiz-questions.csv")
	check(err)
	reader := csv.NewReader(f)
	table, err := reader.ReadAll()
	check(err)
	var questions []question
	for _, row := range table {
		questions = append(questions, question{q: row[0], a: row[1]})
	}
	return questions
}

// ask asks a question and returns an updated score depending on the answer.
func ask(s chan score, question question) {
	fmt.Println("Asking a question...")
	toUpdate := <-s
	fmt.Println(question.q)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter answer: ")
	scanner.Scan()
	text := scanner.Text()
	if strings.Compare(text, question.a) == 0 {
		fmt.Println("Correct!")
		toUpdate++
		s <- toUpdate
	} else {
		fmt.Println("Incorrect :-(")
	}
}

func main() {
	s := score(0)
	qs := questions()
	startTime := time.Now()
	elapsed := 0 * time.Second

	current := 0
	qChannel := make(chan question)
	sChannel := make(chan score)

	// ? Declare to the program that
	// ? sChannel will be listened to
	// ? from here on out
	go ask(sChannel, qs[current])
	sChannel <- s

	finished := false
	for !finished {
		select {
		case question := <-qChannel:
			fmt.Println("Asking question", current+1)
			go ask(sChannel, question)
			current++
			if current == len(qs) {
				finished = true
			} else {
				qChannel <- qs[current]
			}
		case newS := <-sChannel:
			s = newS
		default:
			elapsed = time.Since(startTime)
			if elapsed >= 5*time.Second {
				finished = true
				fmt.Println("You ran out of time!")
			}
		}
	}
	fmt.Println("Final score", s)
}
