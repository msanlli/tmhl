package main

import (
	"time"
	"flag"
	"fmt"
	"os"
	"encoding/csv"
	"bufio"
	"strings"
	"os/signal"
	"syscall"
	"math/rand"
)

type problem struct {
	question string
	answer   string
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Generates random number for the shuffle func
	csvFile, timeLimit := parseFlags()

	// Open and read the csv file
	lines, err := readCSVFile(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file: %v", err))
	}

	// Set and shuffle the qustions
	problems := parseLines(lines)
	problems = shuffle(problems)

	// Play the quiz with the specified questions (csv) and time limit
	playQuiz(problems, *timeLimit)
}

// Flags that enable the user to change the csv file and the time limit
func parseFlags() (*string, *int) {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	return csvFile, timeLimit
}

// csv file reader
func readCSVFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", filename))
	}
	defer file.Close()

	r := csv.NewReader(file)
	return r.ReadAll()
}

// Quiz starter
func playQuiz(problems []problem, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0

	// Handle CTRL+C interruption
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

// Main quiz loop
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		answerCh := make(chan string)
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answerCh <- scanner.Text()
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime's up! You got %d out of %d right.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if strings.TrimSpace(strings.ToLower(answer)) == strings.ToLower(p.answer) {
				correct++
			}
		case <-c:
			fmt.Printf("\nInterrupted! You got %d out of %d right up to this point.\n", correct, i+1)
			break problemLoop
		}
	}

	timer.Stop()
	if correct == len(problems) {
		fmt.Printf("Congratulations! You got all %d questions right. Amazing!\n", len(problems))
	} else {
		fmt.Printf("Quiz ended! You got %d out of %d questions right.\n", correct, len(problems))
	}
}

// Filter the question and the answer on each line of the csv file
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	return ret
}

// Problem shuffling
func shuffle(problems []problem) []problem {
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
	return problems
}

// Handle the fatal errors with an exit message and terminating the program
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
