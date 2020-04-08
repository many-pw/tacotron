package wav

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"math"
)

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
	numChannels := int(f.NumChannels)

	//fmt.Println("len all", len(meta.Data), "/ ba", len(meta.Data)/blockAlign)
	samples = make([]Sample, len(meta.Data)/blockAlign)
	offset := 0
	for i := 0; i < len(samples); i++ {
		soffset := offset + (0 * bitsPerSample / 8)

		var val1 uint
		for b := 0; b*8 < bitsPerSample; b++ {
			val1 += uint(meta.Data[soffset+b]) << uint(b*8)
		}

		samples[i].Values[0] = toInt(val1, bitsPerSample)

		var val2 uint
		if numChannels == 2 {
			soffset = offset + (1 * bitsPerSample / 8)
			for b := 0; b*8 < bitsPerSample; b++ {
				val2 += uint(meta.Data[soffset+b]) << uint(b*8)
			}
		}
		samples[i].Values[1] = toInt(val2, bitsPerSample)
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
