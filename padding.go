package main

//import "fmt"

func pad1D(y []float32, size int, padFlavor string) []float32 {
	s := []float32{}
	s = append(s, y...)
	offset := 2
	j1 := len(y) - offset
	j2 := offset - 1
	for {
		s = append(s, y[j1])
		j1--
		if j1 < 0 {
			j1 += 2
		}
		if len(s) == size+len(y) {
			break
		}
	}
	for {
		s = append([]float32{y[j2]}, s...)
		j2++
		if j2 > len(y)-1 {
			j2 -= 2
		}
		if len(s) == size*2+len(y) {
			break
		}
	}
	return s
}

func pad2d(y [][]float32, size int, padFlavor string) [][]float32 {
	s := [][]float32{}
	j := 0
	size2 := size * 2
	factor := size2 + len(y)
	if len(y)%2 == 0 {
		if size%2 == 0 {
			j = len(y) - 2
		} else {
			j = len(y) - 2
		}
	} else {
		if size%2 == 0 {
			j = len(y) - 1
		} else {
			j = len(y) - 2
		}
	}
	direction := "down"
	special := []float32{}
	for {
		start := append([]float32{}, y[j]...)
		i := 0
		for {
			s0 := float32(start[0])
			start = append(start, s0)
			start = append(start, s0)
			i++
			if size >= i {
				break
			}
		}
		if len(s) == 1 && len(special) == 0 {
			special = append([]float32{}, s[0]...)
			s[0] = append([]float32{}, start...)
		} else {
			//fmt.Println("a", s, "|", start)
			s = append(s, start)
		}
		//fmt.Println(len(s))
		if len(s) >= factor-1 {
			//fmt.Println("a", special, "|", s)
			if len(y)%2 == 0 {
				if size%2 == 0 {
					s = append([][]float32{special}, s...)
				} else {
					s = append(s, special)
				}
			} else {
				if size%2 == 0 {
					s = append([][]float32{special}, s...)
				} else {
					s = append([][]float32{special}, s...)
				}
			}
			break
		}
		if direction == "down" {
			j--
			if j < 0 {
				j += 2
				direction = "up"
			}
		} else {
			j++
			if j > len(y)-1 {
				j -= 2
				direction = "down"
			}
		}
	}
	return s
}
