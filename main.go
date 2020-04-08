package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/guptarohit/asciigraph"
	"github.com/many-pw/tacotron/cycle"
	"github.com/many-pw/tacotron/wav"
)

type Thing struct {
	Count int
	Name  int
}

//var speaker = make(chan float64, 1024*10)
var globalWav = []float64{}
var globalIndex = 0

var global512 = 512
var globalBreak = 0
var globalPause = false
var stream *portaudio.Stream

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	offset, _ := strconv.Atoi(os.Args[2])
	reader := wav.NewReader(file)
	f, meta := reader.Format()
	blocks := float64(len(meta.Data)) / meta.Duration / float64(f.BlockAlign)

	fmt.Println("sr", f.SampleRate, "channels", f.NumChannels)
	fmt.Println("byteRate", f.ByteRate, "BlockAlign", f.BlockAlign)
	fmt.Println("BitsPerSample", f.BitsPerSample)
	samples, _ := reader.ReadSamples(f, meta)
	globalBreak = int(float64(len(samples))/float64(global512)) * 2 // * int(f.BlockAlign) * int(f.NumChannels)
	fmt.Println("duration", meta.Duration, blocks, globalBreak)
	peak := float64(-1.0)
	low := float64(1.0)
	for i, cur := range samples {
		//fmt.Println(cur)
		val := float64(1.0 * reader.FloatValue(f, cur, 0))
		//if rand.Intn(100) > 8 {
		//speaker <- val
		globalWav = append(globalWav, val)
		//speaker <- val
		//}
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
	fmt.Println(offset, peak, low)
	data := []float64{}
	x := 0
	y := len(globalWav) - 1

	x = (len(globalWav) / 2) / 2
	y = len(globalWav) - x - 1

	x = (len(globalWav) / 8) / 2
	y = len(globalWav) - x - 1
	for i, item := range globalWav[x:y] {
		if i%(len(globalWav[x:y])/75) == 0 {
			data = append(data, item)
		}
	}
	graph := asciigraph.Plot(data)

	fmt.Println(graph)
}

func omain() {
	rand.Seed(time.Now().UnixNano())
	portaudio.Initialize()
	if len(os.Args) == 1 {
		fmt.Println("enter 1st param")
		return
	}
	stream, _ = portaudio.OpenDefaultStream(0, 1, 44100, global512, callback)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)
			fmt.Println(command)
			if command == "" {
				globalPause = !globalPause
			} else if command == "-" {
				globalIndex -= globalBreak * 3
				gc -= 3
			}
			if globalPause {
				stream.Stop()
			} else {
				stream.Start()
			}
		}
	}()

	go func() {
		stream.Start()
	}()
	time.Sleep(time.Second * 100)
}

var globalCount = 0
var gc = 0
var globalLast = []float32{}

func process1sec(id int, items []float32) {
	var highsLows = map[int]int{}
	var highsLowsSums = map[int]float32{}
	var highsLowsSum float32
	var dir = ""
	var prevVal = float32(0.0)
	var highLowCount = 0
	for i, val := range items {
		if val > prevVal && dir != "up" {
			dir = "up"
			highLowCount = 0
			highsLowsSums[i] = highsLowsSum
			highsLowsSum = 0
		} else if val < prevVal && dir != "down" {
			dir = "down"
			highLowCount = 0
			highsLowsSums[i] = highsLowsSum
			highsLowsSum = 0
		} else {
			highLowCount += 1
			highsLows[i] = highLowCount
			highsLowsSum += val
		}
		prevVal = val
	}
	things := []Thing{}
	for k, v := range highsLows {
		thing := Thing{v, k}
		things = append(things, thing)
	}
	sort.Slice(things, func(i, j int) bool {
		return things[i].Count > things[j].Count
	})
	/*
		for i, thing := range things {
			fmt.Printf("%v %v %v %0.4f\n", id, thing.Count, thing.Name,
			highsLowsSums[i]/float32(thing.Count))
			if i > 100 {
				break
			}
		}*/
}

func cb(index int) int {
	fmt.Println(index)
	return int(globalWav[globalIndex+index : globalIndex+index+1][0])
}
func callback(_, out []float32) {

	for i, item := range globalWav[globalIndex : globalIndex+global512] {
		out[i] = float32(item)
	}
	a, b := cycle.Brent(cb, 0)
	fmt.Println("Cycle length:", a)
	fmt.Println("Cycle start index:", b)
	globalIndex += global512

	if globalCount*global512 > globalBreak {
		fmt.Println("--- next ---", gc, len(globalLast))
		go process1sec(gc, append([]float32{}, globalLast...))
		globalCount = 0
		globalLast = []float32{}
		gc += 1
	}

	globalCount += 1
	globalLast = append(globalLast, out...)
}
