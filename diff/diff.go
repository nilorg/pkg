package diff

import (
	mapset "github.com/deckarep/golang-set"
)

// IntSlice ...
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

// StringSlice ...
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

// StringSliceV2 ...
func StringSliceV2(src []string, in []string) (added []string, deleted []string, updated []string) {
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
