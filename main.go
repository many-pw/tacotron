package main

import "os"
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
}
