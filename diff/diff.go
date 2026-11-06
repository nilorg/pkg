package diff

import (
	mapset "github.com/deckarep/golang-set"
)

// IntSlice 返回 in 相对 src 的新增项和删除项，自动去重，结果顺序不固定。
func IntSlice(src []int64, in []int64) (added []int64, deleted []int64) {
	srcSet := mapset.NewSet()
	for _, i := range src {
		srcSet.Add(i)
	}
	inSet := mapset.NewSet()
	for _, i := range in {
		inSet.Add(i)
	}
	deletedInterface := srcSet.Difference(inSet).ToSlice()
	for _, i := range deletedInterface {
		deleted = append(deleted, i.(int64))
	}
	addedInterface := inSet.Difference(srcSet).ToSlice()
	for _, i := range addedInterface {
		added = append(added, i.(int64))
	}
	return
}

// IntSliceV2 返回新增项、删除项和交集。
//
// Deprecated: 使用 IntSliceChanges。
func IntSliceV2(src []int64, in []int64) (added []int64, deleted []int64, updated []int64) {
	return IntSliceChanges(src, in)
}

// IntSliceV3 返回非零 ID 的交集和删除项。
//
// Deprecated: 使用 IntSliceUpdatedDeleted。
func IntSliceV3(src []int64, in []int64) (updated []int64, deleted []int64) {
	return IntSliceUpdatedDeleted(src, in)
}

// IntSliceChanges 返回新增项、删除项和交集；updated 不代表字段发生变化。
func IntSliceChanges(src []int64, in []int64) (added []int64, deleted []int64, updated []int64) {
	srcSet := mapset.NewSet()
	for _, i := range src {
		srcSet.Add(i)
	}
	inSet := mapset.NewSet()
	for _, i := range in {
		inSet.Add(i)
	}
	deletedInterface := srcSet.Difference(inSet).ToSlice()
	for _, i := range deletedInterface {
		deleted = append(deleted, i.(int64))
	}
	addedInterface := inSet.Difference(srcSet).ToSlice()
	for _, i := range addedInterface {
		added = append(added, i.(int64))
	}
	updatedInterface := srcSet.Intersect(inSet).ToSlice()
	for _, i := range updatedInterface {
		updated = append(updated, i.(int64))
	}
	return
}

// IntSliceUpdatedDeleted 返回非零 ID 的交集和删除项；updated 不代表字段发生变化。
func IntSliceUpdatedDeleted(src []int64, in []int64) (updated []int64, deleted []int64) {
	srcSet := mapset.NewSet()
	for _, id := range src {
		if id != 0 {
			srcSet.Add(id)
		}
	}
	inSet := mapset.NewSet()
	for _, id := range in {
		if id != 0 {
			inSet.Add(id)
		}
	}
	for _, id := range srcSet.Intersect(inSet).ToSlice() {
		updated = append(updated, id.(int64))
	}
	for _, id := range srcSet.Difference(inSet).ToSlice() {
		deleted = append(deleted, id.(int64))
	}
	return
}

// StringSlice 返回 in 相对 src 的新增项和删除项，自动去重，结果顺序不固定。
func StringSlice(src []string, in []string) (added []string, deleted []string) {
	srcSet := mapset.NewSet()
	for _, i := range src {
		srcSet.Add(i)
	}
	inSet := mapset.NewSet()
	for _, i := range in {
		inSet.Add(i)
	}
	deletedInterface := srcSet.Difference(inSet).ToSlice()
	for _, i := range deletedInterface {
		deleted = append(deleted, i.(string))
	}
	addedInterface := inSet.Difference(srcSet).ToSlice()
	for _, i := range addedInterface {
		added = append(added, i.(string))
	}
	return
}

// StringSliceV2 返回新增项、删除项和交集。
//
// Deprecated: 使用 StringSliceChanges。
func StringSliceV2(src []string, in []string) (added []string, deleted []string, updated []string) {
	return StringSliceChanges(src, in)
}

// StringSliceV3 返回非空 ID 的交集和删除项。
//
// Deprecated: 使用 StringSliceUpdatedDeleted。
func StringSliceV3(src []string, in []string) (updated []string, deleted []string) {
	return StringSliceUpdatedDeleted(src, in)
}

// StringSliceChanges 返回新增项、删除项和交集；updated 不代表字段发生变化。
func StringSliceChanges(src []string, in []string) (added []string, deleted []string, updated []string) {
	srcSet := mapset.NewSet()
	for _, i := range src {
		srcSet.Add(i)
	}
	inSet := mapset.NewSet()
	for _, i := range in {
		inSet.Add(i)
	}
	deletedInterface := srcSet.Difference(inSet).ToSlice()
	for _, i := range deletedInterface {
		deleted = append(deleted, i.(string))
	}
	addedInterface := inSet.Difference(srcSet).ToSlice()
	for _, i := range addedInterface {
		added = append(added, i.(string))
	}
	updatedInterface := srcSet.Intersect(inSet).ToSlice()
	for _, i := range updatedInterface {
		updated = append(updated, i.(string))
	}
	return
}

// StringSliceUpdatedDeleted 返回非空 ID 的交集和删除项；updated 不代表字段发生变化。
func StringSliceUpdatedDeleted(src []string, in []string) (updated []string, deleted []string) {
	srcSet := mapset.NewSet()
	for _, id := range src {
		if id != "" {
			srcSet.Add(id)
		}
	}
	inSet := mapset.NewSet()
	for _, id := range in {
		if id != "" {
			inSet.Add(id)
		}
	}
	for _, id := range srcSet.Intersect(inSet).ToSlice() {
		updated = append(updated, id.(string))
	}
	for _, id := range srcSet.Difference(inSet).ToSlice() {
		deleted = append(deleted, id.(string))
	}
	return
}
