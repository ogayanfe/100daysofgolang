package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type question struct {
	question string
	answer   string
}

// Validates whether the answer to the question is correct or not
func (q *question) checkAnswer(ans string) bool {
	return q.answer == ans
}

var questions []question

// This variable holds a pointer to the csv files to load question from either the one set by user
// or the default given.
var csvFileName = flag.String("csv", "problems.csv", "CSV file to load questions from")

// This variable would hold the quiz time limit.
var maxTime int

// This variable would hold whether to shuffle the questions or not
var shuffle bool

// Reads the answer from Stdin and returns the it
func readAnswer(prefix string) string {
	fmt.Printf("\n%v", prefix) //print question prefix
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func init() {
	const (
		defaultValue = 30
		helpCommand  = "Set quiz time limit in seconds"
	)
	// Set the seed, would be useful for shuffling array
	rand.Seed(time.Now().UnixNano())

	flag.IntVar(&maxTime, "time", defaultValue, helpCommand)
	flag.IntVar(&maxTime, "t", defaultValue, helpCommand)

	flag.BoolVar(&shuffle, "shuffle", false, "Determines whether to shuffle array")
	flag.BoolVar(&shuffle, "s", false, "Shortand for --shuffle")
}

// ask's the user a question, takes in a question type, anschannel to forward result to and
// question number index
func askQuestion(q question, ansChannel chan bool, index int) {
	// ask a question from the user and send answer over channel
	prefix := fmt.Sprintf("Question #%v : %v = ", index+1, q.question)
	ans := readAnswer(prefix)
	ansChannel <- q.checkAnswer(ans)
}

// randomly shuffle array
func shuffleArray(array *[]question) {
	a := *array
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]

	})
}

// Gets quiz questions
func getQuestions() []question {
	questions := loadCSVIntoQuestions(*csvFileName)
	if shuffle {
		shuffleArray(&questions)
	}
	return questions
}

func main() {
	flag.Parse() // Parse CLI flags
	var timer = time.NewTimer(time.Second * time.Duration(maxTime))
	var score int
	var timeUp = false
	ansChannel := make(chan bool)

	questions = getQuestions()
	fmt.Printf("Time Set To %v Second(s)⌛⌛\n", maxTime)

L:
	for i, v := range questions {

		// I'm running function in a goroutine because I want to be able to terminate the program once
		// the time is up even if the user is about to type a new answer. If you don't do this the timer
		// might be up but the user would still be allowed to type in any answer to any question they're
		// currently answering. This way i can interupt user midway if question is wrong
		go askQuestion(v, ansChannel, i)

		// Waits for either a user's time to run out or for an answer to be given
		select {
		case <-timer.C:
			// if the time's up break the loop and continue execution without waiting for the user's response
			// even if they are already typing in their current answer.
			timeUp = true
			break L
		case correct := <-ansChannel:
			if correct {
				score++
			}
		}
	}

	if timeUp {
		fmt.Println("\n\nTime Ran Out😢")
	}
	fmt.Printf("You scored a total of %v out of %v\n", score, len(questions))
}
