package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var apps map[string]int = map[string]int{
	"scikit-lasso":        69,
	"scikit-lda":          53,
	"scikit-ada":          138,
	"scikit-linregr":      45,
	"scikit-rfc":          38,
	"scikit-rfr":          115,
	"spec-astar":          468,
	"spec-leslie":         378,
	"spec-cactus":         780,
	"spec-sphinx":         592,
	"in-memory-analytics": 47,
	"data-serving-client": 72,
	"web-serving-client":  203,
}

var normalized []float64

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	//var counter int
	for scanner.Scan() {
		s := scanner.Text()
		//fmt.Println("Full line\n", s)
		s1 := strings.Split(s, " ")
		appName := s1[0][0 : len(s1[0])-19]
		appTime, _ := strconv.Atoi(s1[1])
		normalized = append(normalized, float64(apps[appName])/float64(appTime))
	}
	//sort the numbers
	sort.Float64s(normalized)
	mNum := len(normalized) / 2
	if isOdd(len(normalized)) {
		fmt.Println(normalized[mNum])
	} else {
		fmt.Println("Median:", (normalized[mNum-1]+normalized[mNum])/2)
	}
}
func isOdd(n int) bool {
	return n%2 == 1
}
