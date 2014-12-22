package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const digits = "0123456789"
const letters = "abcdefghijklmnopqrstvwxyz"

var symbols_set = map[string]string{
	"digits":  digits,
	"letters": letters,
}

var chosen_sets = flag.String("sset", "", "sets of symbols separated by ','")
var symbols = flag.String("s", digits, "symbols")
var max_len = flag.Int("max", 8, "max length")
var min_len = flag.Int("min", 1, "min length")
var file_name = flag.String("f", "", "file to output")

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

func get_progress(cur int, all int) float64 {
	return float64(cur) / float64(all) * 100
}

func output(rounds_count int) {
	if *file_name == "" {
		for i := 0; i < rounds_count; i++ {
			fmt.Println(<-out_chan)
		}
	} else {
		file, err := os.Create(*file_name)
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(file)
		last_print := time.Now().Unix()
		all_bytes := 0
		for i := 0; i < rounds_count; i++ {
			bytes, _ := w.WriteString(<-out_chan + "\n")
			all_bytes = all_bytes + bytes
			if time.Now().Unix() > last_print {
				go func() {
					fmt.Printf("\r%.2f%%, %d MB", get_progress(i, rounds_count), all_bytes/1048576)
				}()
				last_print = time.Now().Unix()
			}
		}
		fmt.Println()
		fmt.Println("done")
		w.Flush()
		file.Close()
	}
	wg.Done()
}

func get_unique(symbols_list []string) []string {
	set := make(map[string]bool)
	for _, val := range symbols_list {
		set[val] = true
	}
	unique := make([]string, 0, len(set))
	for k := range set {
		unique = append(unique, k)
	}
	return unique
}

func get_symbols() string {
	out_string := ""
	if *chosen_sets != "" {
		for _, elem := range strings.Split(*chosen_sets, ",") {
			out_string += symbols_set[elem]
		}
	}
	out_string += *symbols
	return out_string
}

func main() {
	flag.Parse()
	symbols_list = strings.Split(get_symbols(), "")
	symbols_list = get_unique(symbols_list)
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
