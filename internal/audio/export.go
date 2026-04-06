package audio

import (
	"encoding/binary"
	"errors"
	"os"
	"sync"
)

type Recorder struct {
	mu               sync.Mutex
	file             *os.File
	remainingSamples int
	writtenBytes     int
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) Start(path string, durationSeconds float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.file != nil {
		return errors.New("export already in progress")
	}
	if durationSeconds <= 0 {
		return errors.New("duration must be greater than zero")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := writeWAVHeader(file, 0); err != nil {
		_ = file.Close()
		return err
	}

	r.file = file
	r.remainingSamples = int(durationSeconds * SampleRate)
	r.writtenBytes = 0
	return nil
}

func (r *Recorder) WritePCM(data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.file == nil || r.remainingSamples <= 0 {
		return
	}

	maxBytes := r.remainingSamples * 2
	if len(data) > maxBytes {
		data = data[:maxBytes]
	}

	n, err := r.file.Write(data)
	if err != nil {
		r.finishLocked()
		return
	}

	r.writtenBytes += n
	r.remainingSamples -= n / 2
	if r.remainingSamples <= 0 {
		r.finishLocked()
	}
}

func (r *Recorder) finishLocked() {
	if r.file == nil {
		return
	}

	_ = finalizeWAVHeader(r.file, r.writtenBytes)
	_ = r.file.Close()
	r.file = nil
	r.remainingSamples = 0
	r.writtenBytes = 0
}

func writeWAVHeader(file *os.File, dataSize int) error {
	header := make([]byte, 44)
	copy(header[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataSize))
	copy(header[8:12], []byte("WAVE"))
	copy(header[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1)
	binary.LittleEndian.PutUint16(header[22:24], 1)
	binary.LittleEndian.PutUint32(header[24:28], uint32(int(SampleRate)))
	binary.LittleEndian.PutUint32(header[28:32], uint32(int(SampleRate)*2))
	binary.LittleEndian.PutUint16(header[32:34], 2)
	binary.LittleEndian.PutUint16(header[34:36], 16)
	copy(header[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataSize))

	_, err := file.Write(header)
	return err
}

func finalizeWAVHeader(file *os.File, dataSize int) error {
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	return writeWAVHeader(file, dataSize)
}
