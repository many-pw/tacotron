package main

import "github.com/aclements/go-moremath/vec"
import "math"
import "fmt"

func melspectrogram(y []float32) {
	D := stft(y)
	fmt.Println(D)
	//S = amp_to_db(linear_to_mel(np.abs(D)))
	//return normalize(S)
}

func stft(y []float32) string {
	nfft := 2048.0
	w := getWindow(1101)
	fmt.Println(w)
	w = padCenter(w, int(nfft))
	fmt.Println("----")
	fmt.Println(w)
	shape := [][]float64{}
	for _, item := range w {
		shape = append(shape, []float64{item})
	}
	fmt.Println(shape)

	s := [][]float32{}
	s = append(s, []float32{1.0})
	s = append(s, []float32{2.0})
	s = append(s, []float32{3.0})
	s = append(s, []float32{4.0})
	yy := pad(s, 1, "reflect")
	//yy := pad(y, int(math.Ceil(nfft/2.0)), "reflect")
	fmt.Println(yy)
	return "w"
}

func pad(y [][]float32, size int, flavor string) []float32 {
	s := [][]float32{}
	j := len(y) - 2
	direction := "down"
	special := []float32{}
	for {
		start := y[j]
		start = append(start, start[0])
		start = append(start, start[0])
		if len(s) == 1 && len(special) == 0 {
			special = s[0]
			s[0] = start
		} else {
			s = append(s, start)
		}
		fmt.Println(len(s))
		if len(s) >= 5 {
			if len(s)%2 == 0 {
				s = append([][]float32{special}, s...)
			} else {
				s = append(s, special)
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
	fmt.Println(s)
	return []float32{}
}

func padCenter(y []float64, size int) []float64 {
	c := []float64{}
	lpad := int(math.Ceil((float64(size) - float64(len(y))) / 2.0))
	//other := size - len(y) - lpad
	i := 0
	for {
		if len(c) < lpad {
			c = append(c, 0)
		} else {
			if i < len(y) {
				c = append(c, y[i])
				i++
			} else {
				c = append(c, 0)
			}
		}
		if len(c) == size {
			break
		}
	}

	return c
}

func getWindow(length int) []float64 {
	fac := vec.Linspace(math.Pi*-1, math.Pi, length)
	fmt.Println(len(fac))
	w := []float64{}
	alpha := float64(0.5)
	a := [2]float64{alpha, 1.0 - alpha}

	for range fac {
		w = append(w, float64(0.0))
	}
	for k, ak := range a {
		for i, item := range fac {
			val := item * float64(k)
			val = ak * math.Cos(val)
			w[i] = w[i] + val
		}
	}
	return w[0 : len(w)-1]
}
