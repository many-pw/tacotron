package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

func dft(input []float64) []complex128 {
	output := make([]complex128, len(input))

	arg := -2.0 * math.Pi / float64(len(input))
	for k := 0; k < len(input); k++ {
		r, i := 0.0, 0.0
		for n := 0; n < len(input); n++ {
			r += input[n] * math.Cos(arg*float64(n)*float64(k))
			i += input[n] * math.Sin(arg*float64(n)*float64(k))
		}
		output[k] = complex(r, i)
	}
	return output
}

func hfft(samples []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(samples[0], 0)
		return
	}

	//fmt.Printf("%s %d %d\n", "hfft", n, step)
	half := n / 2

	hfft(samples, freqs, half, 2*step)
	hfft(samples[step:], freqs[half:], half, 2*step)

	for k := 0; k < half; k++ {
		a := -2 * math.Pi * float64(k) / float64(n)
		//fmt.Printf("  %s %f\n", "2pi", -2*math.Pi*float64(k)/float64(n))
		e := cmplx.Rect(1, a) * freqs[k+half]
		//fmt.Printf("e %v %v\n", e, cmplx.Rect(1, a))

		freqs[k], freqs[k+half] = freqs[k]+e, freqs[k]-e
	}
}

func fft(samples []float64) []complex128 {
	n := len(samples)
	freqs := make([]complex128, n)
	hfft(samples, freqs, n, 1)
	//for i, cn := range freqs {
	//r, theta := cmplx.Polar(cn)
	//fmt.Printf("%d %.2f %.2f\n", i, r, theta)
	//}
	return freqs
}

func mapRange(n, srcMin, srcMax, dstMin, dstMax float64) float64 {
	return (n-srcMin)/(srcMax-srcMin)*(dstMax-dstMin) + dstMin
}

func drawfft(samples []float64) [][]complex128 {

	items := [][]complex128{}
	bins := 512
	max := 75
	for x := 1; x < max; x++ {
		// n, srcMin, srcMax, dstMin, dstMax
		n0 := int64(mapRange(float64(x-1), 0, float64(max), 0, float64(len(samples))))
		n1 := int64(mapRange(float64(x-0), 0, float64(max), 0, float64(len(samples))))
		// a 0 2968
		c := n0 + (n1-n0)/2
		fmt.Println("a", n0, n1, c)
		sub := make([]float64, bins*2)
		for i := 0; i < len(sub); i += 1 {
			s := 0.0
			n := int(c) - int(bins) + int(i)
			if n >= 0 && n < len(samples) {
				s = samples[n]
			}
			tmp := 1.0
			if true {
				//fmt.Println("2pi-a", float64(i))
				//fmt.Println("2pi- ", float64(i)*math.Pi*2)
				//fmt.Println("2pi- ", float64(i)*math.Pi*2/float64(len(sub)))
				//fmt.Println("2pi- ", math.Cos(float64(i)*math.Pi*2/float64(len(sub))))
				tmp = 0.54 - 0.46*math.Cos(float64(i)*math.Pi*2/float64(len(sub)))
				//fmt.Println("2pi- ", tmp)
			}
			sub[i] = s * tmp
		}

		var freqs []complex128
		if false {
			freqs = dft(sub)
		} else {
			freqs = fft(sub)
		}
		max := 0.0
		for y := 0; y < int(bins); y++ {
			c := freqs[y]
			r2 := cmplx.Abs(c)
			//fmt.Printf("%.2f\n", r2)
			max = math.Max(max, r2)
		}
		//fmt.Printf("max %v %d %.2f\n", int(bins), x, max)
		for y := 0; y < int(bins); y++ {
			c := freqs[y]
			r := 0.0
			if false {
				r = math.Pow(real(c), 2) + math.Pow(imag(c), 2)
			} else {
				r = cmplx.Abs(c)
			}
			if false {
				r = float64(bins) * math.Log10(r/max)
			}
			//img.Set(x, int(bins)-y, gr.ColorAt(r))
			//fmt.Printf(" set %v %d %.2f\n", x, int(bins)-y, r)
		}
	}
	return items
}
