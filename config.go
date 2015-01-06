package main

import (
	"flag"
)

type Config struct {
	chosensets string
	symbols    string
	maxlen     int
	minlen     int
	filename   string
}

const digits = "0123456789"
const letters = "abcdefghijklmnopqrstvwxyz"

var symbols_set = map[string]string{
	"digits":  digits,
	"letters": letters,
}

func (config *Config) parse_cli() {
	chosen_sets_flag := flag.String("sset", "", "sets of symbols separated by ','")
	symbols_flag := flag.String("s", digits, "symbols")
	maxlen_flag := flag.Int("max", 8, "max length")
	minlen_flag := flag.Int("min", 1, "min length")
	filename_flag := flag.String("f", "", "file to output")

	flag.Parse()

	config.chosensets = *chosen_sets_flag
	config.symbols = *symbols_flag
	config.maxlen = *maxlen_flag
	config.minlen = *minlen_flag
	config.filename = *filename_flag
}
