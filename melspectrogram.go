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
	w := getWindow(1101)
	fmt.Println(w)
	w = padCenter(w, 2048)

	return "w"
}

func padCenter(y []float64, size int) []float64 {
	lpad := math.Ceil((float64(size) - float64(len(y))) / 2.0)
	fmt.Println(lpad)

	return []float64{}
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
