package main

import (
	"sort"
	"testing"
)

func slices_compare(first []string, second []string) bool {
	sort.Strings(first)
	sort.Strings(second)
	if len(first) != len(second) {
		return false
	}
	for i, v := range first {
		if second[i] != v {
			return false
		}
	}
	return true
}

func Test_get_symbols(t *testing.T) {
	check := []string{"a", "b", "c", "d", "e", "f"}
	test := get_unique_symbols_list("aaabcdefffddda", "")
	if !slices_compare(check, test) {
		t.Fail()
	}
}

func Test_get_symbols_with_sset(t *testing.T) {
	check := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	test := get_unique_symbols_list("aaabcdefffddda", "digits")
	if !slices_compare(check, test) {
		t.Fail()
	}
}
