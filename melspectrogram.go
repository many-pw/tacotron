package main

import "github.com/aclements/go-moremath/vec"
import "math"
import "fmt"

func melspectrogram(y []float32) {
	D := stft(y)
	fmt.Println(len(D))
	//S = amp_to_db(linear_to_mel(np.abs(D)))
	//return normalize(S)
}

func stft(y []float32) string {
	nfft := 2048.0
	w := getWindow(1101)
	//fmt.Println(w)
	w = padCenter(w, int(nfft))
	//fmt.Println("----")
	//fmt.Println(w)
	shape := [][]float64{}
	for _, item := range w {
		shape = append(shape, []float64{item})
	}
	//fmt.Println(shape)

	//s := []float32{1,2,3}
	//yy := pad1D(s, 6, "reflect")
	//fmt.Println(yy)
	fmt.Println(y[0:3])
	fmt.Println("..")
	yy := pad1D(y, int(math.Ceil(nfft/2.0)), "reflect")
	fmt.Println(yy[0:3])
	fmt.Println(yy[len(yy)-3:])

	s := []float32{1, 2, 3}
	y2d := frame(s, 1, 1)
	fmt.Println(y2d)

	return ""
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
	//fmt.Println(len(fac))
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
