package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	//	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

var allQuestions, correctAnswers int

func askQuestionAndGetAnswer(inLine problem) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(inLine.question)

	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)

	//TODO - trim the input string and do case insensitive checking
	//TODO - how could we validate "almost" correct answers - like in Duolingo, when some typos are allowed?

	//EqualFold - case insensitive compare
	//Trim - remove leading and trailing whitespaces
	if strings.EqualFold(inLine.answer, strings.Trim(text, " ")) {
		correctAnswers++
	}

}

func timeLimitTheQuiz(sec int) {
	timeLimit := time.NewTimer(time.Duration(sec) * time.Second)

	for {
		select {
		case <-timeLimit.C:
			printAndExit(fmt.Sprintf("End of the quiz. You answered %d out of %d questions correctly.", correctAnswers, allQuestions), 1)
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix()) //TODO - more infor about this?

	//Using flags package, create a flag to get path to a quiz file as a parameter. If parameter isn't declared,
	//use the default path that is problems.csv
	quizCSV := flag.String("f", "problems.csv", "Path to a quiz file")
	//Flag to define a custom time of the quiz
	quizTime := flag.Int("t", 30, "Maximum time limit for a quiz")
	//Shuffle the questions
	quizShuffle := flag.Bool("s", false, "Shuffle the order of the questions")
	flag.Parse()

	fmt.Println("Path to selected quiz file is ", *quizCSV)
	fmt.Println("Time limit is", *quizTime, "seconds")
	fmt.Println("Suhuffle ", *quizShuffle)

	var baseOfQ []problem
	//wait for user to press enter to start the timer and show the firs question
	fmt.Printf("Press ENTER to begin the quiz")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		printAndExit(fmt.Sprintf("%s", err), 3)
	}

	//Start the timer thread to exit after the selected period of time
	go timeLimitTheQuiz(*quizTime)

	file, err := os.Open(*quizCSV) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	allQ, err := r.ReadAll()
	if err != nil {
		printAndExit(fmt.Sprintf("%s", err), 2)
	}

	//fill database of questions
	tmpQ := problem{}
	for _, val := range allQ {
		tmpQ.question = val[0]
		tmpQ.answer = val[1]
		baseOfQ = append(baseOfQ, tmpQ)
		allQuestions++
	}

	//TODO - how the hell does the Shuffle parameter function work??
	if *quizShuffle {
		rand.Shuffle(len(baseOfQ), func(i, j int) {
			baseOfQ[i], baseOfQ[j] = baseOfQ[j], baseOfQ[i]
		})
	}

	for _, val := range baseOfQ {
		askQuestionAndGetAnswer(val)
	}
	printAndExit(fmt.Sprintf("End of the quiz. You answered %d out of %d questions correctly.", correctAnswers, allQuestions), 1)
}

func printAndExit(msg string, exitCode int) {
	fmt.Println(msg)
	os.Exit(exitCode)
}
