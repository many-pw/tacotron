package main

import "os"
import "math"
import "fmt"
import "github.com/many-pw/tacotron/wav"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("enter file")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := wav.NewReader(file)
	f, _ := reader.Format()
	fmt.Println("sr", f.SampleRate, "channels", f.NumChannels)
	fmt.Println("byteRate", f.ByteRate, "BlockAlign", f.BlockAlign)
	fmt.Println("BitsPerSample", f.BitsPerSample)
	samples, _ := reader.ReadSamples()
	peak := float64(-1.0)
	low := float64(1.0)
	dir := ""
	prevVal := 0.0
	count := 0
	for _, cur := range samples {
		//fmt.Println(cur)
		val := float64(1.0 * reader.FloatValue(cur, 0))
		if val > prevVal && dir != "up" {
			dir = "up"
			count = 0
		} else if val < prevVal && dir != "down" {
			dir = "down"
			count = 0
		} else {
			count += 1
			fmt.Println("dir is", dir, count)
		}
		prevVal = val
		if val < low {
			low = val
		}
		absVal := math.Abs(float64(val))
		if absVal > peak {
			peak = absVal
		}
	}
	fmt.Println(peak, low)
}
