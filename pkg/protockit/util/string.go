package util

import (
	"strings"
)

func ParseStringSlice(s string) []string {
	items := strings.Split(s, ",")
	for k, v := range items {
		items[k] = strings.Trim(v, " ")
	}
	return Filter(items, func(item string) bool {
		return len(item) > 0
	})
}

func Filter[T any](items []T, fn func(item T) bool) []T {
	res := []T{}
	for _, item := range items {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}
