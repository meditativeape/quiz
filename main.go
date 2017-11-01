package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	var csvPath = flag.String("csv", "problems.csv", "Path to the input CSV file.")
	var timeLimit time.Duration
	flag.Int64Var((*int64)(&timeLimit), "time", 30, "Time limit to solve all quizzes, in seconds.")
	flag.Parse()

	content, err := ioutil.ReadFile(*csvPath)
	check(err)
	inputLines := strings.Split(string(content), "\n")

	totalCount := 0
	correctCount := 0
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Press Enter to start the quiz. You have %d seconds to finish the quiz.\n", timeLimit)
	scanner.Scan()

	go runQuiz(scanner, &totalCount, &correctCount, &inputLines)
	timer := time.NewTimer(time.Second * timeLimit)
	<-timer.C

	fmt.Printf("You correctly answered %d out of %d questions\n", correctCount, totalCount)
}

func runQuiz(scanner *bufio.Scanner, totalCount *int, correctCount *int, inputLines *[]string) {
	for _, quiz := range *inputLines {
		if len(quiz) > 0 {
			(*totalCount)++
			q, a := parseQuiz(quiz)
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
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseQuiz(quiz string) (string, string) {
	parts := strings.Split(quiz, ",")
	return parts[0], parts[1]
}
