package main

import "os"
import "math"
import "fmt"
import "sort"
import "github.com/many-pw/tacotron/wav"

type Thing struct {
	Count int
	Name  int
}

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
	counts := map[int]int{}
	for i, cur := range samples {
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
			//fmt.Printf("dir is %s %d %.4f\n", dir, count, val)
			counts[i] = count
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
	things := []Thing{}
	for k, v := range counts {
		thing := Thing{v, k}
		things = append(things, thing)
	}
	sort.Slice(things, func(i, j int) bool {
		return things[i].Count > things[j].Count
	})
	for i, thing := range things {
		fmt.Println(thing.Count, thing.Name)
		if i > 10000 {
			break
		}
	}
}
