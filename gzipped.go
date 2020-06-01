package gzipped

import (
	"bytes"
	"io"
)

// Scanner knows how to iterate through a multi gzipped file stream, returning a single
// zip file upon each iteration.
type Scanner struct {
	r io.Reader

	fileBuf  *bytes.Buffer // buffers a single file into memory.
	overFlow *bytes.Buffer // contains the start of the next file, if there was an overflow.
	buf      []byte

	bytesRead   int // total bytes read from source.
	byteWritten int // total bytes written.
	matchCount  int // matches consecutive gzip magic bytes.
	size        int // internal buffer size

	err error
}

// NewScanner creates a new *Scanner from a io.Reader.
func NewScanner(r io.Reader, internalBufSize int) *Scanner {
	return &Scanner{
		r:        r,
		fileBuf:  &bytes.Buffer{},
		overFlow: &bytes.Buffer{},
		buf:      make([]byte, internalBufSize),
		size:     internalBufSize,
	}
}

// Scan advances the scanner to the next gzip file, which will then be available through FileBytes()
// if Scan returns true. It returns false if there are no more files to iterate or if the Scanner ran
// into an error. Consult Err() to see what the case is.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	if err := s.loadNextFile(); err != nil {
		return false
	}
	return true
}

// FileBytes returns the buffered file.
func (s *Scanner) FileBytes() []byte {
	return s.fileBuf.Bytes()
}

// Err returns the error which cause the Scanner to fail.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// loadNextFile loads the next zip file into memory, it knows how to stop before beginning in the next zip file.
func (s *Scanner) loadNextFile() error {
	var (
		// gzipMagicBytes are the first 4 bytes of any gzip file.
		gzipMagicBytes = []byte{0x00, 0x00, 0x1f, 0x8b}
		cutoffIdx      = 0
	)

	if err := s.initFileBuf(); err != nil {
		return err
	}

	for {
		n, err := s.read(s.buf)
		if err != nil {
			return err
		}
		for j := 0; j < n; j++ {
			if s.byteWritten == 0 {
				continue
			}
			switch s.matchCount {
			case 0:
				if s.matchByte(s.buf[j], gzipMagicBytes[0]) {
					cutoffIdx = j
				}
				break
			case 1:
				s.matchByte(s.buf[j], gzipMagicBytes[1])
				break
			case 2:
				s.matchByte(s.buf[j], gzipMagicBytes[2])
				break
			case 3:
				if s.matchByte(s.buf[j], gzipMagicBytes[3]) {
					return s.writeCutoff(s.buf, cutoffIdx)
				}
				break
			}
		}
		if err := s.write(s.buf[:n]); err != nil {
			return err
		}
		if s.err == io.EOF {
			return nil
		}
	}
	return nil
}

func (s *Scanner) initFileBuf() error {
	s.fileBuf.Reset()
	if _, err := s.fileBuf.Write(s.overFlow.Bytes()); err != nil {
		return s.setErr(err)
	}
	s.overFlow.Reset()
	return nil
}

func (s *Scanner) read(buf []byte) (int, error) {
	n, err := s.r.Read(buf)
	s.err = err
	s.bytesRead += n
	if err == io.EOF {
		return n, nil
	}
	return n, err
}

func (s *Scanner) write(buf []byte) error {
	n, err := s.fileBuf.Write(buf)
	if err != nil {
		return s.setErr(err)
	}
	s.byteWritten += n
	return nil
}

func (s *Scanner) writeCutoff(buf []byte, cutoff int) error {
	s.matchCount = 0
	if err := s.write(buf[:cutoff]); err != nil {
		return err
	}
	if _, err := s.overFlow.Write(buf[cutoff:]); err != nil {
		return s.setErr(err)
	}
	return nil
}

func (s *Scanner) matchByte(b1, b2 byte) bool {
	if b1 == b2 {
		s.matchCount++
		return true
	}
	s.matchCount = 0
	return false
}

func (s *Scanner) setErr(err error) error {
	s.err = err
	return err
}
