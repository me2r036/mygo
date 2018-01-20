package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	s := make([][]uint8, dx)
	
	for x := range s {
		s[x] = make([]uint8, dy)
		for y := range s[x] {
			s[x][y] = uint8(x * x + y * y)
		}
	}
	
	return s
}

func main() {
	pic.Show(Pic)
}

l