package main

import wave "github.com/zenwerk/go-wave"
import "os"
import "fmt"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("enter file")
		return
	}
	waveReader, _ := wave.NewReader(os.Args[1])
	for {
		framesPerBuffer, _ := waveReader.ReadRawSample()

		fmt.Println(framesPerBuffer[0],
			waveReader.NumSamples, waveReader.ReadSampleNum)
	}
}
