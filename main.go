package main

import "fmt"
import "bufio"
import "strings"
import "io"
import "os"
import "github.com/go-audio/wav"

var wavToWords = map[string]string{}

func main() {
	readFileLines()
	fmt.Println(wavToWords["LJ050-0278"])
	fmt.Println(wavToWords["LJ002-0321"])

	f, err := os.Open("LJ050-0278.wav")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	dur, err := wav.NewDecoder(f).Duration()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s duration: %s\n", f.Name(), dur)

}
func readFileLines() {
	f, err := os.OpenFile("metadata.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

		  fmt.Println(err)
			return
		}
		tokens := strings.Split(line, "|")
		if len(tokens) >= 2 {
			key := tokens[0]
			txt := tokens[len(tokens)-1]
			wavToWords[key] = txt
		}
	}
}
