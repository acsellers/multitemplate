package terse

var Filters = map[string]FilterFunc{}

type FilterFunc func(string) string
