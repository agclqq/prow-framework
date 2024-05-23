package slice

func main() {

}

// Filter 当原切片不会再被使用时，就地 filter 方式是比较推荐的，可以节省内存空间
func Filter() {
	a := []int{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	n := 0

	for _, x := range a {
		if x%2 == 0 {
			a[n] = x
			n++
		}
	}
	a = a[:n]
}
