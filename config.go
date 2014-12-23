package main

import (
	"flag"
)

type Config struct {
	ChosenSets string
	Symbols    string
	MaxLen     int
	MinLen     int
	FileName   string
}

const digits = "0123456789"
const letters = "abcdefghijklmnopqrstvwxyz"

var symbols_set = map[string]string{
	"digits":  digits,
	"letters": letters,
}

func (config *Config) Parse() {
	chosen_sets_flag := flag.String("sset", "", "sets of symbols separated by ','")
	symbols_flag := flag.String("s", digits, "symbols")
	maxlen_flag := flag.Int("max", 8, "max length")
	minlen_flag := flag.Int("min", 1, "min length")
	filename_flag := flag.String("f", "", "file to output")

	flag.Parse()

	config.ChosenSets = *chosen_sets_flag
	config.Symbols = *symbols_flag
	config.MaxLen = *maxlen_flag
	config.MinLen = *minlen_flag
	config.FileName = *filename_flag
}
