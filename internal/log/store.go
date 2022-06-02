package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

/*

# Terminology

The word 'log' can refer to at least three different things: a record, the file that stores
records, and the abstract data type that ties segments together.

To make things less confusing, the following terms are used as defined:

- Record: the data stored in our log.
- Store: the file we store records in.
- Index: the file we store index entries in.
- Segment: the abstraction that ties a store and an index together.
- Log: the abstraction that ties all the segments together.

*/

// enc defines the encoding that we persist record sizes and index entries in.
var enc = binary.BigEndian

// lenWidth defines the number of bytes used to store the record's length.
const (
	lenWidth = 8
)

// store is a simple wrapper around a file with two APIs to append
// and read bytes to and from the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// newStore creates a store for the given file.
func newStore(f *os.File) (*store, error) {
	// Getting the file's information.
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	// Checking the file's current size, in case we're re-creating the store
	// from a file that has existing data, which would happen if, for example,
	// our service had restarted.
	size := uint64(fi.Size())

	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append persists the given bytes to the store.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Storing the start position of the record.
	pos = s.size

	// Write the length of the record so that, when we read the record, we
	// know how many bytes to read.
	//
	// The actual bytes of data are not written but the number of bytes that slice p
	// contains represented as uint64, which is a type that takes up 8 bytes (64 bits)
	// of space. This tells us how many bytes to read, which could be much larger than
	// 8 bytes (uint64 can represent massive numbers), when reading the record.
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	// Write the bytes of data to the buffered writer instead of directly to the file
	// to reduce the number of system calls and improve performance. If a user wrote
	// a lot of small records, this would help a lot.
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}

	// Add the number of bytes used to store the record length (initial binary write)
	// to the number of bytes just written to the buffer. Since uint64 always takes
	// up 8 bytes (64 bits) of space we created a constant that can be referred to
	// called lenWidth. So 'n' accurately represents the total number of bytes
	// written to the buffer.
	n = uint64(w + lenWidth)

	// Updating the store's size to the new number of total bytes written (now represents
	// the start position of the next record).
	s.size = n

	// Return the total number of bytes written (Go APIs convention) and the position
	// where the store holds the record in its file. The segment will use this
	// position when it creates an associated index entry for this record.
	return n, pos, nil
}

// Read returns the record stored at the given position.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Flush method is called to guarantee all data has been forwarded to the
	// underlying io.Writer (our file), in case we're about to try to read a
	// record that the buffer hasn't flushed to disk yet.
	//
	// When we write records, they are initially stored inside an in-memory
	// buffer instead of writing them directly to the disk. This allows for
	// better performance due to not needing to make an operating system call
	// every single time a new record is created. Flush takes the in-memory
	// buffered data and forwards it to the underlying io.Writer to handle
	// which, in our case, is a file belonging to the operating system. This
	// ensures that our persisted records are up to date with the latest writes.
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// Find out how many bytes we have to read to get the whole record. The compiler
	// allocates byte slices that don't escape the functions they're declared in on
	// the stack. A value excapes when it lives beyond the lifetime of the function
	// call - if you return the value, for example.
	//
	// The first 8 bytes of a record contain binary data. That binary data is an
	// encoded uint64 number representing the total amount of bytes stored in the
	// data portion of the record. So we create a slice and read the binary data
	// into it to hold the bytes pertaining to the size.
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	// We take the slice containing the bytes of the the size and convert it to
	// its uint64 representation so we can create a slice large enough to hold
	// the bytes pertaining to the records data. Then we read the data into
	// the slice starting 8 bytes after the given position (so we don't include
	// the bytes representing the size).
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return b, nil
}

// ReadAt reads len(p) bytes into p beginning at the offset in the store's file. It
// implements io.ReaderAt on the store type.
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}

// Close persists any buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}

	return s.File.Close()
}
