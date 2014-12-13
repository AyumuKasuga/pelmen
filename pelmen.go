package main

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"sync"
)

var symbols = flag.String("s", "0123456789", "symbols")
var max_len = flag.Int("max", 8, "max length")
var min_len = flag.Int("min", 1, "min length")

var symbols_list []string
var out_chan chan string = make(chan string, 1000000)

var wg sync.WaitGroup

func slice_to_string(slice []int) string {
	out_string := ""
	for _, elem := range slice {
		out_string += symbols_list[elem]
	}
	return out_string
}

func get_rounds_count(symbols_length int, max_len int) int {
	ret := 0
	for i := 0; i < max_len; i++ {
		ret += int(math.Pow(float64(symbols_length), float64(max_len-i)))
	}
	return ret
}

func output(rounds_count int) {
	for i := 0; i < rounds_count; i++ {
		fmt.Println(<-out_chan)
	}
	wg.Done()
}

func main() {
	flag.Parse()
	symbols_list = strings.Split(*symbols, "")
	symbols_list_length := len(symbols_list)
	max_symbol := symbols_list_length - 1
	index_list := []int{-1}
	rounds_count := get_rounds_count(symbols_list_length, *max_len)

	if *min_len > 1 {
		for i := 1; i < *min_len; i++ {
			index_list = append([]int{0}, index_list...)
		}
		rounds_count = rounds_count - get_rounds_count(symbols_list_length, *min_len-1)
	}

	go output(rounds_count)

	wg.Add(1)

	for i := 0; i < rounds_count; i++ {
		idx := 1
		for {
			symbol_number := index_list[len(index_list)-idx]
			if symbol_number == max_symbol {
				index_list[len(index_list)-idx] = 0
				idx++
				if len(index_list) < idx {
					index_list = append([]int{0}, index_list...)
					break
				}
			} else {
				index_list[len(index_list)-idx]++
				break
			}
		}
		// fmt.Println(index_list)
		out_chan <- slice_to_string(index_list)
		// output(slice_to_string(index_list))
	}
	wg.Wait()
	// fmt.Println("WELL DONE")

}
