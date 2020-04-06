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

var speaker = make(chan float64, 1024)

var global512 = 512
var globalBreak = 0

func main() {
	rand.Seed(time.Now().UnixNano())
	portaudio.Initialize()
	if len(os.Args) == 1 {
		fmt.Println("enter 1st param")
		return
	}
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, global512, callback)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stream)
	if os.Args[1] == "play" {
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := wav.NewReader(file)
	f, meta := reader.Format()
	blocks := float64(len(meta.Data)) / meta.Duration / float64(f.BlockAlign)

	fmt.Println("sr", f.SampleRate, "channels", f.NumChannels)
	fmt.Println("byteRate", f.ByteRate, "BlockAlign", f.BlockAlign)
	fmt.Println("BitsPerSample", f.BitsPerSample)
	samples, _ := reader.ReadSamples(f, meta)
	globalBreak = int(float64(len(samples))/float64(global512)) * int(f.BlockAlign)
	fmt.Println("duration", meta.Duration, blocks, globalBreak)
	peak := float64(-1.0)
	low := float64(1.0)
	dir := ""
	prevVal := 0.0
	count := 0
	counts := map[int]int{}
	go func() {
		stream.Start()
	}()
	for i, cur := range samples {
		//fmt.Println(cur)
		val := float64(4.0 * reader.FloatValue(f, cur, 0))
		//if rand.Intn(100) > 8 {
		speaker <- val
		//speaker <- val
		//}
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
		if i > int(blocks) {
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
	for i, _ := range things {
		//fmt.Println(thing.Count, thing.Name)
		if i > 10 {
			break
		}
	}
	time.Sleep(time.Second * 100)
}

var globalCount = 0
var gc = 0

func callback(_, out []float32) {

	getsome := []float32{}
	for val := range speaker {
		getsome = append(getsome, float32(val))
		if len(getsome) >= global512 {
			break
		}
	}

	if globalCount*global512 > globalBreak {
		fmt.Println("--- next ---", gc)
		globalCount = 0
		gc += 1
	}

	for i := 0; i < global512; i++ {
		out[i] = getsome[i]
	}
	globalCount += 1
}
