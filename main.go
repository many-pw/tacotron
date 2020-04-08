package main

import (
	"bufio"
	"fmt"
	"math/cmplx"
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

	for i := 0; i < len(globalWav); i++ {
		plot(i, globalWav[i])
	}
	waitForSignal()
}

func plot(second int, items []float64) {
	data := []float64{}

	twoD := drawfft(items)
	for _, ca := range twoD {
		data = append(data, cmplx.Abs(ca[0]))
	}

	if len(data) > 0 {
		graph := asciigraph.Plot(data)
		fmt.Println(second)
		fmt.Println(graph)
	}
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

func process1sec(id int, items []float64) {
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
