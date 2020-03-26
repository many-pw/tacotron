package main

import "fmt"
import "bufio"
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
		
		s := stft.New(22050,10)
		complexA := s.STFT(wavSamples)
		fmt.Println(complexA)

	// peak

	/*
	    y = load_wav(path)
    peak = np.abs(y).max()
    if hp.peak_norm or peak > 1.0:
        y /= peak
    mel = melspectrogram(y)
    if hp.voc_mode == 'RAW':
        quant = encode_mu_law(y, mu=2**hp.bits) if hp.mu_law else float_2_label(y, bits=hp.bits)
    elif hp.voc_mode == 'MOL':
        quant = float_2_label(y, bits=16)

    return mel.astype(np.float32), quant.astype(np.int64)

		---- m,x
		 np.save(paths.mel/f'{wav_id}.npy', m, allow_pickle=False)
    np.save(paths.quant/f'{wav_id}.npy', x, allow_pickle=False)
    return wav_id, m.shape[-1]
		*/

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
