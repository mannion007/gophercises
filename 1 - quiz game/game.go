package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"math/rand"
)

var (
	score int = 0
)

type questionAnswer struct {
	question string
	answer   string
}

func main() {

	shuffle := flag.Bool("shuffle", false, "Shuffle the questions?")
	filePath := flag.String("csv", "problems.csv", "Path to the problems file")
	flag.Parse()

	questions := parseQuiz(*filePath, *shuffle)

	reader := bufio.NewReader(os.Stdin)

main:
	for i, q := range questions {
		fmt.Printf("Question %d of %d: %v\n", i+1, len(questions), q.question)

		userInputChan := make(chan string, 1)

		go func() {
			userInput, _ := reader.ReadString('\n')
			userInputChan <- userInput
		}()

		select {
		case <-time.After(time.Duration(10) * time.Second):
			fmt.Printf("You ran out of time!\n")
			break main
		case answer := <-userInputChan:
			if strings.TrimRight(answer, "\n") == q.answer {
				score += 1
			}
		}
	}

	fmt.Printf("You got %d out of a possible %d", score, len(questions))

}

func parseQuiz(filePath string, shuffle bool) []questionAnswer {

	file, err := os.Open(filePath)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	records, err := csv.NewReader(bufio.NewReader(file)).ReadAll()

	if err != nil {
		panic(err)
	}

	parsed := []questionAnswer{}

	for _, rec := range records {
		parsed = append(parsed, questionAnswer{question: rec[0], answer: rec[1]})
	}

	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(parsed), func(i, j int) {
			parsed[i], parsed[j] = parsed[j], parsed[i]
		})
	}

	return parsed
}
