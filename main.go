package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	//	"sort"
	//"strconv"
	"strings"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/guptarohit/asciigraph"
	//	"github.com/many-pw/tacotron/cycle"
	"github.com/many-pw/tacotron/wav"
)

var globalCount = 0
var gc = 0
var globalLast = []float32{}

//var speaker = make(chan float64, 1024*10)
var globalWav = map[int][]float64{}
var globalWavArray = []float64{}
var globalIndex = 0

var global512 = 512
var globalBreak = 0
var globalPause = false
var stream *portaudio.Stream

func main() {
	/* A complex number is a number that can be expressed in the form a + bi, where a and b are real numbers, and i is a solution of the equation x2 = −1.
	 */
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := wav.NewReader(file)
	f, meta := reader.Format()
	samples, _ := reader.ReadSamples(f, meta)

	fmt.Println("")
	fmt.Printf("%20s: %d\n", "Sample Rate", f.SampleRate)
	fmt.Printf("%20s: %d\n", "Samples", len(samples))
	fmt.Printf("%20s: %d\n", "Samples", len(samples)*int(f.BitsPerSample))
	fmt.Printf("%20s: %d\n", "Channels", f.NumChannels)
	fmt.Printf("%20s: %d\n", "BlockAlign", f.BlockAlign)
	fmt.Printf("%20s: %d\n", "ByteRate", f.ByteRate)
	fmt.Printf("%20s: %d %30s\n", "BitRate", f.ByteRate/8, "")
	fmt.Printf("%20s: %d %30s\n", "kb/s", f.ByteRate/125, "")
	fmt.Printf("%20s: %d\n", "BitsPerSample", f.BitsPerSample)
	fmt.Printf("%20s: %f\n", "Duration", meta.Duration)
	fmt.Println("")
	bCounter := 0
	second := 0
	for range samples {
		bCounter += int(f.BitsPerSample)
		if bCounter/1000 >= int(f.ByteRate/125) {
			bCounter = 0
			second += 1
		}
	}
	factor := 1
	if second < int(meta.Duration) {
		factor = 2
	}
	second = 0
	for _, cur := range samples {
		val1 := float64(1.0 * reader.FloatValue(f, cur, 0))
		if f.NumChannels == 2 {
			//val2 := float64(1.0 * reader.FloatValue(f, cur, 1))
		}
		globalWav[second] = append(globalWav[second], val1)
		globalWavArray = append(globalWavArray, val1)

		bCounter += int(f.BitsPerSample)
		if bCounter/1000 >= int(f.ByteRate/125)/factor {
			//fmt.Println(bCounter)
			bCounter = 0
			second += 1
		}
	}
	fmt.Println("remaining", bCounter)
	go startAudio(int(f.SampleRate))

	fmt.Println(len(globalWav))
	for i := 0; i < len(globalWav); i++ {
		plot(i, globalWav[i])
	}
	waitForSignal()
}

func plot(second int, items []float64) {

	data := []float64{}
	x := 0
	y := len(items) - 1

	val := float64(0)
	for i, item := range items[x:y] {
		if i%(len(items[x:y])/75) == 0 {
			data = append(data, val/float64(len(items[x:y])/75))
			val = 0
		} else {
			val += item
		}
	}
	graph := asciigraph.Plot(data)
	fmt.Println(second)
	fmt.Println(graph)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}

func startAudio(sr int) {
	portaudio.Initialize()
	stream, _ = portaudio.OpenDefaultStream(0, 1, float64(sr), global512, callback)
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

	stream.Start()
}

func process1sec(id int, items []float32) {
	/*
			for i, thing := range things {
				fmt.Printf("%v %v %v %0.4f\n", id, thing.Count, thing.Name,
				highsLowsSums[i]/float32(thing.Count))
				if i > 100 {
					break
				}
			}
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
	*/
	data := []float64{}

	for i, item := range items {
		if i%(len(items)/75) == 0 {
			data = append(data, float64(item))
		}
	}
	graph := asciigraph.Plot(data)

	fmt.Println(graph)
	fmt.Println(gc, gc*globalBreak)
}

func callback(_, out []float32) {

	if globalIndex+global512 >= len(globalWavArray) {
		for i := 0; i < len(out); i++ {
			out[i] = 0.0
		}
		return
	}
	for i, item := range globalWavArray[globalIndex : globalIndex+global512] {
		out[i] = float32(item)
	}

	globalIndex += global512
}
