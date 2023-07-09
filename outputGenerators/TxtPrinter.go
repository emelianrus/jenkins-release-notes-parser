package outputGenerators

import (
	"fmt"
	"sort"
)

// returns file with content of plugins from plugin manager

type TxtOutput struct{}

func NewTxtGenerator() TxtOutput {
	return TxtOutput{}
}
func (o TxtOutput) Generate(plugins map[string]string) []byte {

	var data []byte

	// Extract the keys into a slice
	keys := make([]string, 0, len(plugins))
	for key := range plugins {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Iterate over the sorted keys and access the values in the map
	for _, key := range keys {
		value := plugins[key]
		entry := []byte(fmt.Sprintf("%s:%s\n", key, value))
		data = append(data, entry...)

	}
	return data
}
