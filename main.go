package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var csvPath = flag.String("csv", "problems.csv", "Path to the input CSV file.")
	var timeLimit time.Duration
	flag.Int64Var((*int64)(&timeLimit), "time", 30, "Time limit to solve all quizzes, in seconds.")
	var shuffle = flag.Bool("shuffle", false, "Shuffle the order of quizzes.")
	flag.Parse()

	content, err := ioutil.ReadFile(*csvPath)
	check(err)
	inputLines := strings.Split(string(content), "\n")
	quizzes, answers := parseQuiz(&inputLines)

	totalCount := len(*quizzes)
	correctCount := 0
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Press Enter to start the quiz. You have %d seconds to finish the quiz.\n", timeLimit)
	scanner.Scan()

	var wg sync.WaitGroup
	wg.Add(1)
	go runTimer(timeLimit, &wg)
	go runQuiz(scanner, &correctCount, quizzes, answers, *shuffle, &wg)
	wg.Wait() // Wait until either the timer goes off, or user answers all questions

	fmt.Printf("\nYou correctly answered %d out of %d questions\n", correctCount, totalCount)
}

func runTimer(timeLimit time.Duration, wg *sync.WaitGroup) {
	timer := time.NewTimer(time.Second * timeLimit)
	<-timer.C
	(*wg).Done()
}

func parseQuiz(inputLines *[]string) (*[]string, *[]string) {
	quizzes := make([]string, 0)
	answers := make([]string, 0)
	for _, line := range *inputLines {
		if len(line) > 0 {
			parts := strings.Split(line, ",")
			quizzes = append(quizzes, parts[0])
			answers = append(answers, parts[1])
		}
	}
	return &quizzes, &answers
}

func runQuiz(scanner *bufio.Scanner, correctCount *int, quizzes *[]string, answers *[]string, shuffle bool, wg *sync.WaitGroup) {
	rand.Seed(42)
	order := rand.Perm(len(*quizzes))

	for index, value := range order {
		var q, a string
		if shuffle {
			q = (*quizzes)[value]
			a = (*answers)[value]
		} else {
			q = (*quizzes)[index]
			a = (*answers)[index]
		}
		fmt.Printf("Question: %s\n", q)
		fmt.Print("Answer: ")
		scanner.Scan()
		if a == scanner.Text() {
			fmt.Println("Correct!")
			(*correctCount)++
		} else {
			fmt.Println("Wrong...")
		}
	}

	(*wg).Done()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
