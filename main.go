package main

import "github.com/go-audio/wav"
import "os"
import "fmt"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("enter file")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	d := wav.NewDecoder(file)
	dd, e := d.Duration()
	fmt.Println(d.Format(), dd, e)
}
