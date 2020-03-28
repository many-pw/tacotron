package wav

import (
	"io"
)

const (
	AudioFormatPCM       = 1
	AudioFormatIEEEFloat = 3
)

type WavFormat struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type WavData struct {
	io.Reader
	Size uint32
	pos  uint32
}

type Sample struct {
	Values [2]int
}

type RIFFReader interface {
	io.Reader
	io.ReaderAt
}

type TheRiffReader struct {
	RIFFReader
}

type RIFFChunk struct {
	FileSize uint32
	FileType []byte
	Chunks   []*Chunk
}

type Chunk struct {
	ChunkID   []byte
	ChunkSize uint32
	RIFFReader
}
