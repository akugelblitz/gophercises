package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	filePtr := flag.String("filename", "problems.csv", "Name of the questions file. Default problems.csv")
	timePtr := flag.Int("time", 30, "Time limit of quiz in seconds. Default 30.")
	flag.Parse()

	fmt.Println("Time limit", *timePtr, "seconds.Press any key to continue.")
	score, totalQues := quizzerWithTimeout(*filePtr, *timePtr)
	fmt.Println("final score:", score, "out of", totalQues)
}

func quizzerWithTimeout(fileName string, timeLimit int) (int, int) {
	questions := getQuestions(fileName)
	totalQues := len(questions)

	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeLimit))

	defer cancel()
	go quizzer(ctx, ch, questions)

	var score int = 0
	
	for{
		breakFlag := false
		select {
		case <-ctx.Done():
			fmt.Printf("Context cancelled: %v\n", ctx.Err())
			// for i := range ch {
			// 	fmt.Println(i)
			// 	score += i
			// }
			breakFlag = true
	
		case res, ok := <-ch:
			if !ok {
				// return score, totalQues
				breakFlag = true
			}
			score += res
		}
		if breakFlag{
			break
		}
	}

	return score, totalQues
}

func quizzer(ctx context.Context, ch chan int, questions [][]string) {
	var userAns string
	for i := 0; i < len(questions); i++ {
		fmt.Println(questions[i][0])
		fmt.Scanln(&userAns)
		ch <- checkAns(userAns, questions[i][1])
	}
	close(ch)
}

func getQuestions(fileName string) [][]string {
	fd, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened text file")
	defer fd.Close()

	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	return records
}

func checkAns(userAns string, correctAns string) int {
	if userAns == correctAns {
		return 1
	}
	return 0
}
