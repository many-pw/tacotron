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

func (r *Reader) ReadSamples(f *WavFormat, meta *WavMeta) (samples []Sample, err error) {
	blockAlign := int(f.BlockAlign)
	bitsPerSample := int(f.BitsPerSample)

	fmt.Println("len all", len(meta.Data), "/ ba", len(meta.Data)/blockAlign)
	samples = make([]Sample, len(meta.Data)/blockAlign)
	offset := 0
	j := 0
	for i := 0; i < len(samples); i++ {
		soffset := offset + (j * bitsPerSample / 8)

		var val uint
		for b := 0; b*8 < bitsPerSample; b++ {
			//fmt.Println(soffset + b)
			val += uint(meta.Data[soffset+b]) << uint(b*8)
		}

		samples[i].Values[j] = toInt(val, bitsPerSample)
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

func (r *Reader) Format() (*WavFormat, *WavMeta) {
	f, _ := r.readFormat()
	sampleRate := int(f.SampleRate)
	bitsPerSample := int(f.BitsPerSample)
	numChannels := int(f.NumChannels)
	data, _ := r.readData()
	meta := WavMeta{}
	meta.Data, _ = ioutil.ReadAll(data)
	meta.Duration = float64(len(meta.Data)) / float64(sampleRate*numChannels*bitsPerSample/8)
	return f, &meta
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

	//fmt.Println(r.WavData.Size, len(p))

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

func (r *Reader) FloatValue(f *WavFormat, sample Sample, channel uint) float32 {
	return float32(float64(r.IntValue(sample, channel)) / math.Pow(2, float64(f.BitsPerSample)))
}
