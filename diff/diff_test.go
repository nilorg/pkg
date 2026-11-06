package diff

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func assertInts(t *testing.T, name string, got, want []int64) {
	t.Helper()
	actual := append([]int64(nil), got...)
	expected := append([]int64(nil), want...)
	sort.Slice(actual, func(i, j int) bool { return actual[i] < actual[j] })
	sort.Slice(expected, func(i, j int) bool { return expected[i] < expected[j] })
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s: got %v, want %v", name, actual, expected)
	}
}

func TestIntSliceFunctions(t *testing.T) {
	tests := []struct {
		name                             string
		src, in, added, deleted, updated []int64
		nonzeroUpdated, nonzeroDeleted   []int64
	}{
		{name: "nil slices"},
		{name: "empty slices", src: []int64{}, in: []int64{}},
		{name: "add all", in: []int64{2, 1}, added: []int64{1, 2}},
		{name: "delete all", src: []int64{2, 1}, deleted: []int64{1, 2}, nonzeroDeleted: []int64{1, 2}},
		{name: "identical", src: []int64{2, 1}, in: []int64{1, 2}, updated: []int64{1, 2}, nonzeroUpdated: []int64{1, 2}},
		{name: "disjoint", src: []int64{1}, in: []int64{2}, added: []int64{2}, deleted: []int64{1}, nonzeroDeleted: []int64{1}},
		{name: "mixed duplicates", src: []int64{1, 1, 2, 3}, in: []int64{2, 3, 4, 4}, added: []int64{4}, deleted: []int64{1}, updated: []int64{2, 3}, nonzeroUpdated: []int64{2, 3}, nonzeroDeleted: []int64{1}},
		{name: "zero intersection", src: []int64{0, 0, 1}, in: []int64{0, 2, 0}, added: []int64{2}, deleted: []int64{1}, updated: []int64{0}, nonzeroDeleted: []int64{1}},
		{name: "zero deleted", src: []int64{0, 1}, in: []int64{1}, deleted: []int64{0}, updated: []int64{1}, nonzeroUpdated: []int64{1}},
		{name: "zero added", src: []int64{1}, in: []int64{0, 1}, added: []int64{0}, updated: []int64{1}, nonzeroUpdated: []int64{1}},
		{name: "zero only", src: []int64{0}, in: []int64{0}, updated: []int64{0}},
		{name: "negative and limits", src: []int64{math.MinInt64, -1}, in: []int64{-1, math.MaxInt64}, added: []int64{math.MaxInt64}, deleted: []int64{math.MinInt64}, updated: []int64{-1}, nonzeroUpdated: []int64{-1}, nonzeroDeleted: []int64{math.MinInt64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcBefore := append([]int64{}, tt.src...)
			inBefore := append([]int64{}, tt.in...)
			added, deleted := IntSlice(tt.src, tt.in)
			assertInts(t, "added", added, tt.added)
			assertInts(t, "deleted", deleted, tt.deleted)
			genericAdded, genericDeleted := Slice(tt.src, tt.in)
			assertInts(t, "generic added", genericAdded, tt.added)
			assertInts(t, "generic deleted", genericDeleted, tt.deleted)
			for _, fn := range []func([]int64, []int64) ([]int64, []int64, []int64){IntSliceV2, IntSliceChanges, SliceChanges[int64]} {
				added, deleted, updated := fn(tt.src, tt.in)
				assertInts(t, "added", added, tt.added)
				assertInts(t, "deleted", deleted, tt.deleted)
				assertInts(t, "updated", updated, tt.updated)
			}
			for _, fn := range []func([]int64, []int64) ([]int64, []int64){IntSliceV3, IntSliceUpdatedDeleted, SliceUpdatedDeleted[int64]} {
				updated, deleted := fn(tt.src, tt.in)
				assertInts(t, "updated", updated, tt.nonzeroUpdated)
				assertInts(t, "deleted", deleted, tt.nonzeroDeleted)
			}
			if !reflect.DeepEqual(append([]int64{}, tt.src...), srcBefore) || !reflect.DeepEqual(append([]int64{}, tt.in...), inBefore) {
				t.Fatal("function modified its input")
			}
		})
	}
}

func TestGenericSliceChanges(t *testing.T) {
	type ID uint64
	added, deleted, updated := SliceChanges([]ID{1, 2, 2}, []ID{2, 3, 3})
	if !reflect.DeepEqual(added, []ID{3}) || !reflect.DeepEqual(deleted, []ID{1}) || !reflect.DeepEqual(updated, []ID{2}) {
		t.Fatalf("unexpected changes: %v %v %v", added, deleted, updated)
	}
}

func TestGenericEmptyResultsAreNil(t *testing.T) {
	added, deleted := Slice[int](nil, nil)
	if added != nil || deleted != nil {
		t.Fatal("Slice should return nil for empty results")
	}
	added, deleted, updated := SliceChanges[int](nil, nil)
	if added != nil || deleted != nil || updated != nil {
		t.Fatal("SliceChanges should return nil for empty results")
	}
	updated, deleted = SliceUpdatedDeleted([]int{0}, []int{0})
	if updated != nil || deleted != nil {
		t.Fatal("SliceUpdatedDeleted should return nil for empty results")
	}
}

func assertStrings(t *testing.T, name string, got, want []string) {
	t.Helper()
	actual := append([]string(nil), got...)
	expected := append([]string(nil), want...)
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%s: got %q, want %q", name, actual, expected)
	}
}

func TestStringSliceFunctions(t *testing.T) {
	tests := []struct {
		name                             string
		src, in, added, deleted, updated []string
		nonemptyUpdated, nonemptyDeleted []string
	}{
		{name: "nil slices"},
		{name: "empty slices", src: []string{}, in: []string{}},
		{name: "add all", in: []string{"b", "a"}, added: []string{"a", "b"}},
		{name: "delete all", src: []string{"b", "a"}, deleted: []string{"a", "b"}, nonemptyDeleted: []string{"a", "b"}},
		{name: "identical", src: []string{"b", "a"}, in: []string{"a", "b"}, updated: []string{"a", "b"}, nonemptyUpdated: []string{"a", "b"}},
		{name: "disjoint", src: []string{"a"}, in: []string{"b"}, added: []string{"b"}, deleted: []string{"a"}, nonemptyDeleted: []string{"a"}},
		{name: "mixed changes", src: []string{"a", "b", "c"}, in: []string{"c", "d", "b"}, added: []string{"d"}, deleted: []string{"a"}, updated: []string{"b", "c"}, nonemptyUpdated: []string{"b", "c"}, nonemptyDeleted: []string{"a"}},
		{name: "duplicates", src: []string{"a", "a", "b", "b"}, in: []string{"a", "a", "c", "c"}, added: []string{"c"}, deleted: []string{"b"}, updated: []string{"a"}, nonemptyUpdated: []string{"a"}, nonemptyDeleted: []string{"b"}},
		{name: "empty intersection", src: []string{"", "", "a"}, in: []string{"", "b", ""}, added: []string{"b"}, deleted: []string{"a"}, updated: []string{""}, nonemptyDeleted: []string{"a"}},
		{name: "empty deleted", src: []string{"", "a"}, in: []string{"a"}, deleted: []string{""}, updated: []string{"a"}, nonemptyUpdated: []string{"a"}},
		{name: "empty added", src: []string{"a"}, in: []string{"", "a"}, added: []string{""}, updated: []string{"a"}, nonemptyUpdated: []string{"a"}},
		{name: "empty only", src: []string{""}, in: []string{""}, updated: []string{""}},
		{name: "exact string matching", src: []string{"01", "A", " ", "中文"}, in: []string{"1", "a", " ", "中文"}, added: []string{"1", "a"}, deleted: []string{"01", "A"}, updated: []string{" ", "中文"}, nonemptyUpdated: []string{" ", "中文"}, nonemptyDeleted: []string{"01", "A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions := []struct {
				name string
				run  func(*testing.T)
			}{
				{"StringSliceV2", func(t *testing.T) {
					added, deleted, updated := StringSliceV2(tt.src, tt.in)
					assertStrings(t, "added", added, tt.added)
					assertStrings(t, "deleted", deleted, tt.deleted)
					assertStrings(t, "updated", updated, tt.updated)
				}},
				{"StringSliceV3", func(t *testing.T) {
					updated, deleted := StringSliceV3(tt.src, tt.in)
					assertStrings(t, "updated", updated, tt.nonemptyUpdated)
					assertStrings(t, "deleted", deleted, tt.nonemptyDeleted)
				}},
				{"StringSlice", func(t *testing.T) {
					added, deleted := StringSlice(tt.src, tt.in)
					assertStrings(t, "added", added, tt.added)
					assertStrings(t, "deleted", deleted, tt.deleted)
				}},
				{"StringSliceChanges", func(t *testing.T) {
					added, deleted, updated := StringSliceChanges(tt.src, tt.in)
					assertStrings(t, "added", added, tt.added)
					assertStrings(t, "deleted", deleted, tt.deleted)
					assertStrings(t, "updated", updated, tt.updated)
				}},
				{"StringSliceUpdatedDeleted", func(t *testing.T) {
					updated, deleted := StringSliceUpdatedDeleted(tt.src, tt.in)
					assertStrings(t, "updated", updated, tt.nonemptyUpdated)
					assertStrings(t, "deleted", deleted, tt.nonemptyDeleted)
				}},
			}
			for _, fn := range functions {
				t.Run(fn.name, func(t *testing.T) {
					srcBefore := append([]string{}, tt.src...)
					inBefore := append([]string{}, tt.in...)
					fn.run(t)
					if !reflect.DeepEqual(append([]string{}, tt.src...), srcBefore) || !reflect.DeepEqual(append([]string{}, tt.in...), inBefore) {
						t.Fatal("function modified its input")
					}
				})
			}
		})
	}
}
