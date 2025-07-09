package output

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"
)

type BufferedJSONLWriter struct {
	mu     sync.Mutex
	writer *bufio.Writer
	closer io.Closer
}

func NewBufferedJSONLWriter(w io.WriteCloser, size int) *BufferedJSONLWriter {
	return &BufferedJSONLWriter{
		writer: bufio.NewWriterSize(w, size),
		closer: w,
	}
}

func (b *BufferedJSONLWriter) WriteSnapshot(snap Snapshot) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	data, err := json.Marshal(snap)
	if err != nil {
		return err
	}

	_, err = b.writer.Write(data)
	if err != nil {
		return err
	}

	return b.writer.WriteByte('\n')
}

func (b *BufferedJSONLWriter) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.writer.Flush()
}

func (b *BufferedJSONLWriter) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	_ = b.writer.Flush()
	return b.closer.Close()
}
