package diff

import mapset "github.com/deckarep/golang-set/v3"

// Slice 返回 in 相对 src 的新增项和删除项，自动去重，结果顺序不固定。
func Slice[T comparable](src, in []T) (added, deleted []T) {
	srcSet := mapset.NewSet(src...)
	inSet := mapset.NewSet(in...)
	return srcDifference(inSet, srcSet), srcDifference(srcSet, inSet)
}

// SliceChanges 返回新增项、删除项和交集，自动去重，结果顺序不固定。
// updated 不代表字段发生变化。
func SliceChanges[T comparable](src, in []T) (added, deleted, updated []T) {
	srcSet := mapset.NewSet(src...)
	inSet := mapset.NewSet(in...)
	return srcDifference(inSet, srcSet), srcDifference(srcSet, inSet), nilIfEmpty(srcSet.Intersect(inSet).ToSlice())
}

// SliceUpdatedDeleted 返回非零值的交集和删除项，自动去重，结果顺序不固定。
// updated 不代表字段发生变化。
func SliceUpdatedDeleted[T comparable](src, in []T) (updated, deleted []T) {
	srcSet := mapset.NewSet(src...)
	inSet := mapset.NewSet(in...)
	var zero T
	srcSet.Remove(zero)
	inSet.Remove(zero)
	return nilIfEmpty(srcSet.Intersect(inSet).ToSlice()), srcDifference(srcSet, inSet)
}

func srcDifference[T comparable](src, in mapset.Set[T]) []T {
	return nilIfEmpty(src.Difference(in).ToSlice())
}

// nilIfEmpty 保持旧接口在没有结果时返回 nil 的行为。
func nilIfEmpty[T any](values []T) []T {
	if len(values) == 0 {
		return nil
	}
	return values
}

// IntSlice 返回 in 相对 src 的新增项和删除项，自动去重，结果顺序不固定。
func IntSlice(src, in []int64) (added, deleted []int64) {
	return Slice(src, in)
}

// IntSliceV2 返回新增项、删除项和交集。
//
// Deprecated: 使用 IntSliceChanges。
func IntSliceV2(src, in []int64) (added, deleted, updated []int64) {
	return IntSliceChanges(src, in)
}

// IntSliceV3 返回非零 ID 的交集和删除项。
//
// Deprecated: 使用 IntSliceUpdatedDeleted。
func IntSliceV3(src, in []int64) (updated, deleted []int64) {
	return IntSliceUpdatedDeleted(src, in)
}

// IntSliceChanges 返回新增项、删除项和交集；updated 不代表字段发生变化。
func IntSliceChanges(src, in []int64) (added, deleted, updated []int64) {
	return SliceChanges(src, in)
}

// IntSliceUpdatedDeleted 返回非零 ID 的交集和删除项；updated 不代表字段发生变化。
func IntSliceUpdatedDeleted(src, in []int64) (updated, deleted []int64) {
	return SliceUpdatedDeleted(src, in)
}

// StringSlice 返回 in 相对 src 的新增项和删除项，自动去重，结果顺序不固定。
func StringSlice(src, in []string) (added, deleted []string) {
	return Slice(src, in)
}

// StringSliceV2 返回新增项、删除项和交集。
//
// Deprecated: 使用 StringSliceChanges。
func StringSliceV2(src, in []string) (added, deleted, updated []string) {
	return StringSliceChanges(src, in)
}

// StringSliceV3 返回非空 ID 的交集和删除项。
//
// Deprecated: 使用 StringSliceUpdatedDeleted。
func StringSliceV3(src, in []string) (updated, deleted []string) {
	return StringSliceUpdatedDeleted(src, in)
}

// StringSliceChanges 返回新增项、删除项和交集；updated 不代表字段发生变化。
func StringSliceChanges(src, in []string) (added, deleted, updated []string) {
	return SliceChanges(src, in)
}

// StringSliceUpdatedDeleted 返回非空 ID 的交集和删除项；updated 不代表字段发生变化。
func StringSliceUpdatedDeleted(src, in []string) (updated, deleted []string) {
	return SliceUpdatedDeleted(src, in)
}
