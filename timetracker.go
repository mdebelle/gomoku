package main

import (
	"time"
	"fmt"
	"text/tabwriter"
	"os"
)

type funcTime struct {
	totalTime	time.Duration
	count		uint
}

var funcTimes map[string]funcTime = make(map[string]funcTime)

func timeFunc(start time.Time, funcName string) {
	elapsed := time.Since(start)
	if _, ok := funcTimes[funcName]; ok {
		tmp := funcTimes[funcName]
		funcTimes[funcName] = funcTime{tmp.totalTime + elapsed, tmp.count + 1}
	} else {
		funcTimes[funcName] = funcTime{elapsed, 1}
	}
}

func resetTimer() {
	funcTimes = make(map[string]funcTime)
}

func displayAverages() {
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 8, 0, 2, ' ', 0)
	fmt.Println("---------- Function times ----------")
	
	for key, timer := range funcTimes {
		fmt.Fprintln(w, key, "\t", timer.count, "\t", timer.totalTime, "\t", timer.totalTime / time.Duration(timer.count))
	}
	w.Flush()
	fmt.Println("------------------------------------")
}
