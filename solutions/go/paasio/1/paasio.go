package paasio

import (
	"io"
	"sync"
)

type readCounter struct {
	reader io.Reader
	count  int64
	ops    int
	m      sync.Mutex
}

type writeCounter struct {
	writer io.Writer
	count  int64
	ops    int
	m      sync.Mutex
}

type readWriteCounter struct {
	reader *readCounter
	writer *writeCounter
}

func (rwc *readWriteCounter) Read(p []byte) (int, error) {
	return rwc.reader.Read(p)
}

func (rwc *readWriteCounter) Write(p []byte) (int, error) {
	return rwc.writer.Write(p)
}

func (rwc *readWriteCounter) ReadCount() (int64, int) {
	return rwc.reader.ReadCount()
}

func (rwc *readWriteCounter) WriteCount() (int64, int) {
	return rwc.writer.WriteCount()
}

func NewWriteCounter(writer io.Writer) WriteCounter {
	return &writeCounter{writer, 0, 0, sync.Mutex{}}
}

func NewReadCounter(reader io.Reader) ReadCounter {
	return &readCounter{reader, 0, 0, sync.Mutex{}}
}

func NewReadWriteCounter(readwriter io.ReadWriter) ReadWriteCounter {
	return &readWriteCounter{
		NewReadCounter(readwriter).(*readCounter),
		NewWriteCounter(readwriter).(*writeCounter)}
}

func (rc *readCounter) Read(p []byte) (int, error) {
	nr, err := rc.reader.Read(p)
	rc.m.Lock()
	defer rc.m.Unlock()
	rc.count += int64(nr)
	rc.ops++
	return nr, err
}

func (rc *readCounter) ReadCount() (int64, int) {
	rc.m.Lock()
	defer rc.m.Unlock()
	return rc.count, rc.ops
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	nr, err := wc.writer.Write(p)
	wc.m.Lock()
	defer wc.m.Unlock()
	wc.count += int64(nr)
	wc.ops++
	return nr, err
}

func (wc *writeCounter) WriteCount() (int64, int) {
	return wc.count, wc.ops
}
