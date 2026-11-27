package sorter

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"

	"siguma0013/reskk-dictionary/internal/model"
)

// SortEntries : ソート済辞書データスライスを生成
func SortEntries[T model.Entry](entries []T) []T {
	entries = slices.Clone(entries)

	slices.SortStableFunc(entries, func(left, right T) int {
		return cmp.Compare(left.GetIndexYomi(), right.GetIndexYomi())
	})

	return entries
}

// IsSorted : 辞書データがソート済か判定
func IsSorted[T model.Entry](entries []T) []error {
	sortedEntries := SortEntries(entries)

	var differences []error

	for i, entry := range entries {
		if !reflect.DeepEqual(entry, sortedEntries[i]) {
			differences = append(differences, fmt.Errorf("  - [%d] %s（%s） -> %s（%s）", i, entry.GetKanji(), entry.GetIndexYomi(), sortedEntries[i].GetKanji(), sortedEntries[i].GetIndexYomi()))
		}
	}

	return differences
}
