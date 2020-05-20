package wav

import (
	"encoding/binary"
	"io"
)

type writeCallback func(w io.Writer)

func NewTheRiffWriter(w io.Writer, fileType []byte, fileSize uint32) *RiffWriter {
	w.Write([]byte("RIFF"))
	binary.Write(w, binary.LittleEndian, fileSize)
	w.Write(fileType)

	return &RiffWriter{w}
}

func (w *RiffWriter) WriteChunk(chunkID []byte, chunkSize uint32, cb writeCallback) (err error) {
	_, err = w.Write(chunkID)

	if err != nil {
		return
	}

	err = binary.Write(w, binary.LittleEndian, chunkSize)

	if err != nil {
		return
	}

	cb(w)

	return
}
