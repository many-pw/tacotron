package main

import "os"
import "math"
import "fmt"
import "time"
import "sort"
import "math/rand"
import "github.com/many-pw/tacotron/wav"
import "github.com/gordonklaus/portaudio"

type Thing struct {
	Count int
	Name  int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	portaudio.Initialize()
	if len(os.Args) == 1 {
		fmt.Println("enter 1st param")
		return
	}
	if os.Args[1] == "play" {
		stream, err := portaudio.OpenDefaultStream(0, 1, 44100, 512, callback)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(stream)

		go func() {
			stream.Start()
		}()

		for {
			time.Sleep(1)
		}
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := wav.NewReader(file)
	f, meta := reader.Format()
	fmt.Println("duration", meta.Duration)

	fmt.Println("sr", f.SampleRate, "channels", f.NumChannels)
	fmt.Println("byteRate", f.ByteRate, "BlockAlign", f.BlockAlign)
	fmt.Println("BitsPerSample", f.BitsPerSample)
	samples, _ := reader.ReadSamples(f, meta)
	peak := float64(-1.0)
	low := float64(1.0)
	dir := ""
	prevVal := 0.0
	count := 0
	counts := map[int]int{}
	for i, cur := range samples {
		//fmt.Println(cur)
		val := float64(1.0 * reader.FloatValue(f, cur, 0))
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

var global int

func callback(_, out []float32) {
	fmt.Println(1, len(out))
	for i := 0; i < 512; i++ {
		if global > 10 {
			out[i] = 0.99
		} else {
			out[i] = 0
		}
		fmt.Println(out[i])
	}
	if global == 11 {
		global = 0
	}
	global++
}
