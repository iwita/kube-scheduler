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
	"scikit-lasso":        55,
	"scikit-lda":          39,
	"scikit-ada":          69,
	"scikit-linregr":      33,
	"scikit-rfc":          36,
	"scikit-rfr":          100,
	"spec-astar":          345,
	"spec-leslie":         299,
	"spec-cactus":         479,
	"spec-sphinx":         466,
	"in-memory-analytics": 47,
	"data-serving-client": 38,
	"web-serving-client":  203,
	"spec-mcf":            214,
	"spec-lbm":            349,
}

var normalized []float64

func main() {
	offset, _ := strconv.Atoi(os.Args[1])
	scanner := bufio.NewScanner(os.Stdin)
	//var counter int
	for scanner.Scan() {
		s := scanner.Text()
		//fmt.Println("Full line\n", s)
		s1 := strings.Split(s, " ")
		appName := s1[0][0 : len(s1[0])-offset]
		appTime, _ := strconv.Atoi(s1[1])
		normalized = append(normalized, float64(apps[appName])/float64(appTime))
		fmt.Printf("%v %v %v %v\n", appName, float64(apps[appName])/float64(appTime), s1[2], s1[3])

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
