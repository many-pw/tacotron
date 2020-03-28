package wav

import "errors"
import "bufio"
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

	bytes = make([]byte, numSamples*blockAlign)
	fmt.Println("aaa", len(bytes))
	n, err = r.Read(bytes)

	if err != nil {
		return
	}

	numSamples = n / blockAlign
	r.WavData.pos += uint32(numSamples * blockAlign)
	samples = make([]Sample, numSamples)
	offset := 0

	for i := 0; i < numSamples; i++ {
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

				var val uint
				for b = 0; b*8 < bitsPerSample; b++ {
					val += uint(bytes[soffset+b]) << uint(b*8)
				}

				samples[i].Values[j] = toInt(val, bitsPerSample)
			}
		}

		offset += blockAlign
	}

	return
}

func findChunk(riffChunk *RIFFChunk, id string) (chunk *Chunk) {
	for i, ch := range riffChunk.Chunks {
		fmt.Println(i, string(ch.ChunkID[:]))
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
