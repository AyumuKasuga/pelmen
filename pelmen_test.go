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

func Test_get_unique(t *testing.T) {
	unique := get_unique([]string{"a", "a", "c", "b", "b", "d", "a", "b", "c", "c"})
	if !slices_compare(unique, []string{"a", "b", "c", "d"}) {
		t.Fail()
	}
}
