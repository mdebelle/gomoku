package main

import (
	"time"
	"fmt"
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

func displayAverages() {
	for key, timer := range funcTimes {
		fmt.Println(key, " : ", timer.totalTime / time.Duration(timer.count))
	}
}
