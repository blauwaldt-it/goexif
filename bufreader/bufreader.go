package bufreader

import (
	"errors"
	"io"
)

type ReadSeekReaderAt interface {
	io.ReadSeeker
	io.ReaderAt
}

type cachedReader interface {
	Seek(offset int64, whence int) (int64, error)
	Read(p []byte, off int64) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)
	Size() int64
	IncClients()
	DecClients()
}

type seqReader struct {
	r         ReadSeekReaderAt
	fSize     int64
	chunkSize int

	buffer  []byte
	clients int
}

var ErrOutOfRange = errors.New("Bufreader: seek out of range")
var ErrReadErr = errors.New("Bufreader: error reading")

//////////////////////////////////////////////////////////////////////////////

func NewSeqReaderSize(r ReadSeekReaderAt, chunkSize int) (sr *seqReader) {
	fSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		fSize = 0
	}

	sr = &seqReader{
		r:         r,
		fSize:     fSize,
		chunkSize: chunkSize,
		clients:   1,
	}

	return
}

func NewSeqReader(r ReadSeekReaderAt) (sr *seqReader) {
	return NewSeqReaderSize(r, 16*1024)
}

func (rs *seqReader) Seek(offset int64, whence int) (int64, error) {
	return rs.r.Seek(offset, whence)
}

func (rs *seqReader) Read(p []byte, off int64) (n int, err error) {
	if rs.fSize <= 0 {
		return 0, ErrReadErr
	}
	if off > rs.fSize {
		return 0, ErrOutOfRange
	}

	rs.r.Seek(off, io.SeekStart)
	return rs.r.Read(p)
}

func (rs *seqReader) ReadAt(p []byte, off int64) (n int, err error) {
	if rs.fSize <= 0 {
		return 0, ErrReadErr
	}
	if off > rs.fSize {
		return 0, ErrOutOfRange
	}

	rs.r.Seek(off, io.SeekStart)
	return rs.r.Read(p)
}

func (rs *seqReader) Size() int64 {
	return rs.fSize
}

func (rs *seqReader) IncClients() {
	rs.clients++
}

func (rs *seqReader) DecClients() {
	if rs.clients >= 1 {
		rs.clients--
		if rs.clients == 0 {
			rs.fSize = 0
			closer, ok := rs.r.(io.Closer)
			if ok {
				closer.Close()
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////////

type Bufreader struct {
	cr cachedReader

	offset int64 // Start of this reader absolute in cr
	seek   int64 // Always from offset
}

func NewSize(r ReadSeekReaderAt, bufSize int) (bufrd *Bufreader) {
	newbufrd := &Bufreader{
		cr: NewSeqReaderSize(r, bufSize),
	}

	bufrd = newbufrd
	return
}

func New(r ReadSeekReaderAt) (bufrd *Bufreader) {
	return NewSize(r, 16*1024)
}

func (bufrd *Bufreader) Seek(offset int64, whence int) (int64, error) {
	return bufrd.cr.Seek(offset, whence)
}

func (bufrd *Bufreader) Read(p []byte) (n int, err error) {
	n, err = bufrd.cr.Read(p, bufrd.seek+bufrd.offset)
	bufrd.seek += int64(n)
	return
}

func (bufrd *Bufreader) ReadAt(p []byte, off int64) (n int, err error) {
	return bufrd.cr.ReadAt(p, bufrd.seek+bufrd.offset)
}

func (bufrd *Bufreader) Close() {
	bufrd.cr.DecClients()
}

func (bufrd *Bufreader) NewReaderOffsAbs(offs int64) (newBufrd *Bufreader, err error) {
	if offs < 0 || offs > bufrd.cr.Size() {
		return nil, ErrOutOfRange
	}

	bufrd.cr.IncClients()
	newBufrd = &Bufreader{
		cr:     bufrd.cr,
		offset: offs,
	}
	return
}

func (bufrd *Bufreader) NewReaderOffsRel(offsd int64) (newBufrd *Bufreader, err error) {
	return bufrd.NewReaderOffsAbs(bufrd.seek + offsd)
}

func (bufrd *Bufreader) NewReaderOffs() (newBufrd *Bufreader) {
	newBufrd, _ = bufrd.NewReaderOffsAbs(bufrd.seek)
	return
}
