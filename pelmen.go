package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// var symbols_list []string
// var out_chan chan string = make(chan string, 1000000)

// var wg sync.WaitGroup

// func slice_to_string(slice []int) string {
// 	out_string := ""
// 	for _, elem := range slice {
// 		out_string += symbols_list[elem]
// 	}
// 	return out_string
// }

// func get_rounds_count(symbols_length int, max_len int) int {
// 	ret := 0
// 	for i := 0; i < max_len; i++ {
// 		ret += int(math.Pow(float64(symbols_length), float64(max_len-i)))
// 	}
// 	return ret
// }

func get_progress(cur int, all int) float64 {
	return float64(cur) / float64(all) * 100
}

func output(rounds_count int, file_name string, wg *sync.WaitGroup, out_chan chan string) {
	if file_name == "" {
		for i := 0; i < rounds_count; i++ {
			fmt.Println(<-out_chan)
		}
	} else {
		file, err := os.Create(file_name)
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

// func get_unique_symbols_list(symbols string, chosen_sets string) []string {
// 	out_string := ""
// 	if chosen_sets != "" {
// 		for _, elem := range strings.Split(chosen_sets, ",") {
// 			out_string += symbols_set[elem]
// 		}
// 	}
// 	out_string += symbols
// 	symbols_list = strings.Split(out_string, "")

// 	set := make(map[string]bool)
// 	for _, val := range symbols_list {
// 		set[val] = true
// 	}
// 	unique := make([]string, 0, len(set))
// 	for k := range set {
// 		unique = append(unique, k)
// 	}

// 	return unique
// }

func main() {
	config := Config{}
	config.parse_cli()
	gen := Generator{config: config}
	gen.init()
	var wg sync.WaitGroup
	wg.Add(1)
	var out_chan chan string = make(chan string, 1000000)
	go output(gen.rounds_count, gen.config.filename, &wg, out_chan)
	for i := 0; i < gen.rounds_count; i++ {
		gen.next()
		out_chan <- gen.get_string()
	}
	wg.Wait()
	// symbols_list = get_unique_symbols_list(config.Symbols, config.ChosenSets)
	// symbols_list_length := len(symbols_list)
	// max_symbol := symbols_list_length - 1
	// index_list := []int{-1}
	// rounds_count := get_rounds_count(symbols_list_length, config.MaxLen)

	// if config.MinLen > 1 {
	// 	for i := 1; i < config.MinLen; i++ {
	// 		index_list = append([]int{0}, index_list...)
	// 	}
	// 	rounds_count = rounds_count - get_rounds_count(symbols_list_length, config.MinLen-1)
	// }

	// go output(rounds_count, config.FileName)

	// wg.Add(1)

	// for i := 0; i < rounds_count; i++ {
	// 	idx := 1
	// 	for {
	// 		symbol_number := index_list[len(index_list)-idx]
	// 		if symbol_number == max_symbol {
	// 			index_list[len(index_list)-idx] = 0
	// 			idx++
	// 			if len(index_list) < idx {
	// 				index_list = append([]int{0}, index_list...)
	// 				break
	// 			}
	// 		} else {
	// 			index_list[len(index_list)-idx]++
	// 			break
	// 		}
	// 	}
	// 	// fmt.Println(index_list)
	// 	out_chan <- slice_to_string(index_list)
	// 	// output(slice_to_string(index_list))
	// }
	// wg.Wait()
	// // fmt.Println("WELL DONE")

}

type Generator struct {
	config              Config
	symbols_list        []string
	symbols_list_length int
	max_symbol          int
	index_list          []int
	rounds_count        int
	round               int
}

func (gen *Generator) init() {
	gen.populate_symbols_list()
	gen.max_symbol = gen.symbols_list_length - 1
	gen.calculate_rounds_count()
	gen.index_list = []int{-1}
	gen.round = 1
}

func (gen *Generator) populate_symbols_list() {
	out_string := ""
	if gen.config.chosensets != "" {
		for _, elem := range strings.Split(gen.config.chosensets, ",") {
			out_string += symbols_set[elem]
		}
	}
	out_string += gen.config.symbols
	gen.symbols_list = strings.Split(out_string, "")

	set := make(map[string]bool)
	for _, val := range gen.symbols_list {
		set[val] = true
	}
	unique := make([]string, 0, len(set))
	for k := range set {
		unique = append(unique, k)
	}

	gen.symbols_list = unique
	gen.symbols_list_length = len(gen.symbols_list)
}

func (gen *Generator) get_rounds_count(symbols_length int, max_len int) int {
	ret := 0
	for i := 0; i < max_len; i++ {
		ret += int(math.Pow(float64(symbols_length), float64(max_len-i)))
	}
	return ret
}

func (gen *Generator) calculate_rounds_count() {
	count := gen.get_rounds_count(gen.symbols_list_length, gen.config.maxlen)
	if gen.config.minlen > 1 {
		for i := 1; i < gen.config.minlen; i++ {
			gen.index_list = append([]int{0}, gen.index_list...)
		}
		count = count - gen.get_rounds_count(gen.symbols_list_length, gen.config.minlen-1)
	}
	gen.rounds_count = count
}

func (gen *Generator) next() {
	idx := 1
	for {
		symbol_number := gen.index_list[len(gen.index_list)-idx]
		if symbol_number == gen.max_symbol {
			gen.index_list[len(gen.index_list)-idx] = 0
			idx++
			if len(gen.index_list) < idx {
				gen.index_list = append([]int{0}, gen.index_list...)
				break
			}
		} else {
			gen.index_list[len(gen.index_list)-idx]++
			break
		}
	}
}

func (gen *Generator) get_string() string {
	out_string := ""
	for _, elem := range gen.index_list {
		out_string += gen.symbols_list[elem]
	}
	return out_string
}
