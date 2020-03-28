package wav

import "errors"
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

func (r *Reader) readFormat() (fmt *WavFormat, err error) {
	var riffChunk *RIFFChunk

	fmt = new(WavFormat)

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

	err = binary.Read(fmtChunk, binary.LittleEndian, fmt)
	if err != nil {
		return
	}

	return
}
