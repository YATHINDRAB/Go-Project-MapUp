// main.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type requestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type responsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortSequential(toSort [][]int) [][]int {
	for i := range toSort {
		sort.Ints(toSort[i])
	}
	return toSort
}

func sortConcurrent(toSort [][]int) [][]int {
	var wg sync.WaitGroup

	for i := range toSort {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sort.Ints(toSort[i])
		}(i)
	}

	wg.Wait()
	return toSort
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var reqPayload requestPayload
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		fmt.Println("Error decoding request payload:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
	return 
	}

	startTime := time.Now()
	sortedArrays := sortSequential(reqPayload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	respPayload := responsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respPayload)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var reqPayload requestPayload
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortConcurrent(reqPayload.ToSort)
	timeTaken := time.Since(startTime).Nanoseconds()

	respPayload := responsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respPayload)
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	port := ":8000"
	http.ListenAndServe(port, nil)
}

