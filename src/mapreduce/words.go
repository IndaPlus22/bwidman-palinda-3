package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"
const Routines = 32

// Return the word frequencies of the text argument.
//
// Split load optimally across processor cores.
func WordCount(text string) map[string]int {
	words := strings.Fields(text)
	var freqMaps [Routines]map[string]int
	var wg sync.WaitGroup
	
	// Map
	sectionLen := len(words) / Routines
	for i := 0; i < Routines-1; i++ {
		wg.Add(1)
		go mapping(words[i*sectionLen:(i+1)*sectionLen], &freqMaps[i], &wg)
	}
	wg.Add(1)
	go mapping(words[(Routines-1)*sectionLen:], &freqMaps[Routines-1], &wg) // Count rest at the end
	wg.Wait()
	
	// Reduce
	reduceCh := make(chan map[string]int)
	go reduce(freqMaps[:len(freqMaps) / 2], freqMaps[len(freqMaps) / 2:], reduceCh)
	freqs := <-reduceCh

	return freqs
}

func mapping(section []string, freqMap *map[string]int, wg *sync.WaitGroup) {
	*freqMap = make(map[string]int)
	for _, word := range section {
		trimmed := strings.Trim(strings.ToLower(word), ".,")
		(*freqMap)[trimmed]++
	}
	wg.Done()
}

func reduce(freqMapsLeft []map[string]int, freqMapsRight []map[string]int, reduceCh chan<- map[string]int) {
	m := len(freqMapsLeft)
	n := len(freqMapsRight)
	if m < 2 && n < 2 {
		reduceCh <- merge(freqMapsLeft[0], freqMapsRight[0])
		return
	}
	mergedCh := make(chan map[string]int)
	// Recursively split down to individual maps
	go reduce(freqMapsLeft[:m/2], freqMapsLeft[m/2:], mergedCh)
	go reduce(freqMapsRight[:n/2], freqMapsRight[n/2:], mergedCh)
	// Return the merged map
	reduceCh <- merge(<-mergedCh, <-mergedCh)
}

func merge(freqMap1 map[string]int, freqMap2 map[string]int) (mergedMap map[string]int) {
	mergedMap = freqMap1
	for word, freq := range freqMap2 {
		mergedMap[word] += freq
	}
	return
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	bytes, err := os.ReadFile(DataFile)
	if err != nil {
		return
	}
	data := string(bytes)

	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}