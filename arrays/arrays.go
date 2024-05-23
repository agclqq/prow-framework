package arrays

import "reflect"

// InArray 判断obj是否在数组/map target中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func DistinctInt(s []int) []int {
	set := make(map[int]bool)
	res := make([]int, 0)
	for _, v := range s {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = true
		res = append(res, v)
	}
	return res
}

func DistinctStr(s []string) []string {
	set := make(map[string]bool)
	res := make([]string, 0)
	for _, v := range s {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = true
		res = append(res, v)
	}
	return res
}

func Intersect(arr1, arr2 []string) []string {
	m := make(map[string]bool)
	rs := make([]string, 0)

	for _, v := range arr1 {
		m[v] = true
	}
	for _, v := range arr2 {
		if _, ok := m[v]; ok {
			rs = append(rs, v)
		}
	}
	return rs
}

// SumUint uint列表求和
func SumUint(slice []uint) uint {
	var sum uint
	for _, value := range slice {
		sum += value
	}
	return sum
}
