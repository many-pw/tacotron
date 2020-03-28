package wav

import "errors"
import "bufio"
import "io/ioutil"
import "fmt"
import "math"
import "encoding/binary"

type Reader struct {
	rr        *TheRiffReader
	riffChunk *RIFFChunk
	format    *WavFormat
	*WavData
}

func NewReader(r RIFFReader) *Reader {
	riffReader := NewTheRiffReader(r)
	return &Reader{rr: riffReader}
}

func (r *Reader) ReadSamples(params ...uint32) (samples []Sample, err error) {
	format, err := r.Format()
	if err != nil {
		return
	}

	//numSamples := 22050
	blockAlign := int(format.BlockAlign)
	bitsPerSample := int(format.BitsPerSample)

	fmt.Printf("blockAlign: %d\n", blockAlign)
	fmt.Printf("bitsPerSample: %d\n", bitsPerSample)

	data, _ := r.readData()
	fmt.Println(data)
	all, _ := ioutil.ReadAll(data)
	fmt.Println(len(all))
	samples = make([]Sample, len(all)/2)
	offset := 0
	j := 0
	for i := 0; i < len(samples); i++ {
		soffset := offset + (j * bitsPerSample / 8)

		var val uint
		for b := 0; b*8 < bitsPerSample; b++ {
			val += uint(all[soffset+b]) << uint(b*8)
		}

		samples[i].Values[j] = toInt(val, bitsPerSample)
		offset += blockAlign
	}

	return
}
func (r *Reader) OldReadSamples(params ...uint32) (samples []Sample, err error) {
	var bytes []byte
	var numSamples, b, n int

	if len(params) > 0 {
		numSamples = int(params[0])
	} else {
		numSamples = 2048
	}

	format, err := r.Format()
	if err != nil {
		return
	}

	numChannels := int(format.NumChannels)
	blockAlign := int(format.BlockAlign)
	bitsPerSample := int(format.BitsPerSample)

	fmt.Printf("numChannels: %d\n", numChannels)
	fmt.Printf("blockAlign: %d\n", blockAlign)
	fmt.Printf("bitsPerSample: %d\n", bitsPerSample)

	bytes = make([]byte, numSamples*blockAlign)
	n, err = r.Read(bytes)

	if err != nil {
		return
	}

	numSamples = n / blockAlign
	fmt.Printf("i read %d bytes.\n", n)
	fmt.Printf("%d / 2 = %d \n", n, n/2)
	fmt.Printf("pos: %d\n", r.WavData.pos)
	r.WavData.pos += uint32(numSamples * blockAlign)
	fmt.Printf("pos: %d\n", r.WavData.pos)
	samples = make([]Sample, numSamples)
	offset := 0

	for i := 0; i < numSamples; i++ {
		//fmt.Printf("i: %d\n", i)
		if format.AudioFormat == AudioFormatIEEEFloat {
			for j := 0; j < numChannels; j++ {
				soffset := offset + (j * bitsPerSample / 8)

				bits :=
					uint32((int(bytes[soffset+3]) << 24) +
						(int(bytes[soffset+2]) << 16) +
						(int(bytes[soffset+1]) << 8) +
						int(bytes[soffset]))
				samples[i].Values[j] = int(math.MaxInt32 * math.Float32frombits(bits))
			}
		} else {
			for j := 0; j < numChannels; j++ {
				soffset := offset + (j * bitsPerSample / 8)
				//fmt.Printf("offset, j, j * bitsPerSample, f: %v, %v, %v, %v\n",
				//offset, j, bitsPerSample/8, j*bitsPerSample/8)

				var val uint
				for b = 0; b*8 < bitsPerSample; b++ {
					val += uint(bytes[soffset+b]) << uint(b*8)
					//fmt.Printf("b: %v %v\n", b, val)
				}

				samples[i].Values[j] = toInt(val, bitsPerSample)
			}
		}

		offset += blockAlign
	}

	return
}

func findChunk(riffChunk *RIFFChunk, id string) (chunk *Chunk) {
	for _, ch := range riffChunk.Chunks {
		if string(ch.ChunkID[:]) == id {
			chunk = ch
			break
		}
	}

	return
}

func (r *Reader) Format() (format *WavFormat, err error) {
	if r.format == nil {
		format, err = r.readFormat()
		if err != nil {
			return
		}
		r.format = format
	} else {
		format = r.format
	}

	return
}

func (r *Reader) readFormat() (wfmt *WavFormat, err error) {
	var riffChunk *RIFFChunk

	wfmt = new(WavFormat)

	if r.riffChunk == nil {
		riffChunk, err = r.rr.Read()
		if err != nil {
			return
		}

		r.riffChunk = riffChunk
	} else {
		riffChunk = r.riffChunk
	}

	fmtChunk := findChunk(riffChunk, "fmt ")

	if fmtChunk == nil {
		err = errors.New("Format chunk is not found")
		return
	}

	err = binary.Read(fmtChunk, binary.LittleEndian, wfmt)
	if err != nil {
		return
	}

	return
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.WavData == nil {
		data, err := r.readData()
		if err != nil {
			return n, err
		}
		r.WavData = data
	}

	fmt.Println(r.WavData.Size, len(p))

	return r.WavData.Read(p)
}

func (r *Reader) readData() (data *WavData, err error) {
	var riffChunk *RIFFChunk

	if r.riffChunk == nil {
		riffChunk, err = r.rr.Read()
		if err != nil {
			return
		}

		r.riffChunk = riffChunk
	} else {
		riffChunk = r.riffChunk
	}

	dataChunk := findChunk(riffChunk, "data")
	if dataChunk == nil {
		err = errors.New("Data chunk is not found")
		return
	}

	data = &WavData{bufio.NewReader(dataChunk), dataChunk.ChunkSize, 0}

	return
}

func (r *Reader) IntValue(sample Sample, channel uint) int {
	return sample.Values[channel]
}

func (r *Reader) FloatValue(sample Sample, channel uint) float32 {
	return float32(float64(r.IntValue(sample, channel)) / math.Pow(2, float64(r.format.BitsPerSample)))
}
