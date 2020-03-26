package main

import "fmt"
import "bufio"
import "github.com/aclements/go-moremath/vec"
import "math/cmplx"
import "math"
import "strings"
import "io"
import "os"
import "github.com/youpy/go-wav"
import "github.com/r9y9/gossp/stft"
//import "github.com/go-audio/wav"
// for mat.Dot import "gonum.org/v1/gonum/mat"

var wavToWords = map[string]string{}

func main() {
	readFileLines()
	fmt.Println(wavToWords["LJ050-0278"])
	fmt.Println(wavToWords["LJ002-0321"])

	/*
	f, err := os.Open("LJ002-0321.wav")
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
*/

    file, err := os.Open("LJ002-0321.wav")
    if err != nil {
      panic(err)
    }
    reader := wav.NewReader(file)
    wavformat, err_rd := reader.Format()
    if err_rd != nil {
        panic(err_rd)
    }

    if wavformat.AudioFormat != wav.AudioFormatPCM {
        panic("Audio format is invalid ")
    }

    fmt.Println("Block align is ", wavformat.BlockAlign)

    samples, err := reader.ReadSamples(22050) // 2048
    wavSamples := make([]float64, 0)

    for _, curr_sample := range samples {
       wavSamples = append(wavSamples, reader.FloatValue(curr_sample, 0))
    }
		
		ns := uint16(len(wavSamples))/wavformat.NumChannels
		fmt.Println(wavformat.BitsPerSample, wavformat.SampleRate, wavformat.NumChannels)


		peak := float64(0.0)
		for _, val := range wavSamples {
			a := math.Abs(val)
			if a > peak {
				peak = a
			}
		}
		fmt.Println(peak, ns)
		
		nfft := 20
		s := stft.New(22050,nfft)
		complexD := []float64{}
		for _, listOfComplex := range s.STFT(wavSamples) {
			for _, val := range listOfComplex {
			  complexD = append(complexD, cmplx.Abs(val))
			}
		}
		fmt.Println(complexD)
		complexS := []float64{}
		for _,listOfComplex := range s.STFT(complexD) {
			for _, val := range listOfComplex {
			  complexS = append(complexS, cmplx.Abs(val))
			}
		}
		fmt.Println(complexS)

		nmels := 20
		//  weights = np.zeros((n_mels, int(1 + n_fft // 2)), dtype=dtype)
		a := int(1.0+(float64(nfft) / 2.0))
		fmt.Println(nmels, a)

		a2d := [][]float64{}
		for i := 0; i<nmels; i++ {
			thing := []float64{}
			for j :=0; j<a; j++ {
				thing = append(thing, 0.0)
			}
			a2d = append(a2d, thing)
		}

		fmt.Println(a2d)

		//each one should be a length
		//therer should be nmels of them
	  
		fftfreqs := vec.Linspace(0.0, 11025.0, 11) 
		fmt.Println(fftfreqs)

		min_mel := hz_to_mel(0)
		max_mel := hz_to_mel(11025.0)

		fmt.Println(min_mel, max_mel)

	/*
	    y = load_wav(path)
    peak = np.abs(y).max()
    if hp.peak_norm or peak > 1.0:
        y /= peak
    mel = melspectrogram(y)
		*/

}

func hz_to_mel(a float64) float64 {
	f_min := 0.0
	f_sp := 200.0 / 3.0

	mels := (a - f_min) / f_sp

	min_log_hz := 1000.0    
	min_log_mel := (min_log_hz - f_min) / f_sp 
	logstep := math.Log(6.4) / 27.0

	mels = min_log_mel + math.Log(a / min_log_hz) / logstep

	return mels
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
