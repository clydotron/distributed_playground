package main

import "fmt"

func FindOdd(seq []int) int {

	m := make(map[int]int)

	for _, s := range seq {

		v, ok := m[s]
		fmt.Printf("adding %d, ok=%v\n", s, ok)
		if ok {
			m[s] = v + 1
		} else {
			m[s] = 1
		}
	}
	fmt.Println(m)
	for k, v := range m {
		if v%2 == 1 {

			return k
		}
	}
	return 0 // your code here
}

func main() {
	s := []int{20, 1, -1, 2, -2, 3, 3, 5, 5, 1, 2, 4, 20, 4, -1, -2, 5}
	FindOdd(s)
}
