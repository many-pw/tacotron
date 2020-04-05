package main

import wave "github.com/zenwerk/go-wave"
import "os"
import "math"
import "fmt"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("enter file")
		return
	}
	waveReader, _ := wave.NewReader(os.Args[1])
	bitsPerSample := 16
	for {
		framesPerBuffer, _ := waveReader.ReadRawSample()

		for i, f := range framesPerBuffer {

			val := float64(f) / math.Pow(2, float64(bitsPerSample))

			fmt.Println(val, i, f,
				waveReader.NumSamples, waveReader.ReadSampleNum)
		}
	}
}
