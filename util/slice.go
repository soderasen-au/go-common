package util

func SumSlice(s []int) int {
	sum := 0
	for _, i := range s {
		sum = sum + i
	}
	return sum
}
