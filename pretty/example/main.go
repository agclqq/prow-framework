package main

import (
	"fmt"
	"time"

	"github.com/agclqq/prow-framework/pretty"
)

func main() {
	progressBar()

}

func progressBar() {
	fmt.Println("\nthis is a progress bar inline")
	pretty.ProgressBar(50, 100*time.Millisecond)
	time.Sleep(1 * time.Second)
	fmt.Println("\nthis is a progress bar with new line")
	pretty.ProgressBarWithNewLine(50, 100*time.Millisecond)

	fmt.Println("\nthis is a progress bar with clear screen")
	time.Sleep(1 * time.Second)
	pretty.ProgressBarWithClearScreen(50, 100*time.Millisecond)

	time.Sleep(1 * time.Second)
	fmt.Println("\nthis is a progress bar rate")
	for i := 0.1; i <= 1; i += 0.1 {
		pretty.ProgressBarRate(50, i)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("\nthis is a progress bar rate with new line")
	for i := 0.1; i <= 1; i += 0.1 {
		pretty.ProgressBarRateWithNewLine(50, i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\nthis is a progress bar rate with clear screen")
	time.Sleep(1 * time.Second)
	for i := 0.1; i <= 1; i += 0.1 {
		pretty.ProgressBarRateWithClearScreen(50, i)
		time.Sleep(1 * time.Second)
	}

}
