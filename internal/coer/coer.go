// Package coer implements the small subset of the Canonical Octet Encoding
// Rules (ITU-T X.696) needed for the IEEE 1609.2 structures in this project.
package coer

import (
	"encoding/binary"
	"errors"
)

// ErrTruncated is returned when the input ends before a value is complete.
var ErrTruncated = errors.New("coer: truncated input")

// ErrFormat is returned for encodings this subset does not accept.
var ErrFormat = errors.New("coer: malformed encoding")

// Writer appends COER-encoded values to B.
type Writer struct{ B []byte }

func (w *Writer) U8(v uint8)   { w.B = append(w.B, v) }
func (w *Writer) U16(v uint16) { w.B = binary.BigEndian.AppendUint16(w.B, v) }
func (w *Writer) U32(v uint32) { w.B = binary.BigEndian.AppendUint32(w.B, v) }
func (w *Writer) U64(v uint64) { w.B = binary.BigEndian.AppendUint64(w.B, v) }

// Raw appends bytes without a length (fixed-size OCTET STRING).
func (w *Writer) Raw(b []byte) { w.B = append(w.B, b...) }

// Len appends a length determinant (X.696 clause 8.6).
func (w *Writer) Len(n int) {
	if n < 0x80 {
		w.B = append(w.B, byte(n))
		return
	}
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], uint64(n))
	i := 0
	for tmp[i] == 0 {
		i++
	}
	w.B = append(w.B, 0x80|byte(8-i))
	w.B = append(w.B, tmp[i:]...)
}

// Var appends a variable-size OCTET STRING or an open type: length
// determinant followed by the bytes.
func (w *Writer) Var(b []byte) {
	w.Len(len(b))
	w.B = append(w.B, b...)
}

// UintUnbounded appends an INTEGER (0..MAX): length determinant plus the
// minimal big-endian octets (X.696 clause 10.4 d).
func (w *Writer) UintUnbounded(v uint64) {
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], v)
	i := 0
	for i < 7 && tmp[i] == 0 {
		i++
	}
	w.Len(8 - i)
	w.B = append(w.B, tmp[i:]...)
}

// Choice appends a CHOICE tag for alternative index idx (context-specific
// class, idx < 63).
func (w *Writer) Choice(idx int) { w.B = append(w.B, 0x80|byte(idx)) }

// Count appends the quantity field of a SEQUENCE OF.
func (w *Writer) Count(n int) { w.UintUnbounded(uint64(n)) }

// Reader decodes COER values from a byte slice. The first error sticks.
type Reader struct {
	b   []byte
	off int
	err error
}

func NewReader(b []byte) *Reader { return &Reader{b: b} }

func (r *Reader) Err() error  { return r.err }
func (r *Reader) Offset() int { return r.off }
func (r *Reader) Left() int   { return len(r.b) - r.off }

func (r *Reader) take(n int) []byte {
	if r.err != nil {
		return nil
	}
	if n < 0 || r.off+n > len(r.b) {
		r.err = ErrTruncated
		return nil
	}
	s := r.b[r.off : r.off+n]
	r.off += n
	return s
}

func (r *Reader) fail() { // mark malformed
	if r.err == nil {
		r.err = ErrFormat
	}
}

// Fail records a format error raised by the caller.
func (r *Reader) Fail() { r.fail() }

func (r *Reader) U8() uint8 {
	s := r.take(1)
	if s == nil {
		return 0
	}
	return s[0]
}

func (r *Reader) U16() uint16 {
	s := r.take(2)
	if s == nil {
		return 0
	}
	return binary.BigEndian.Uint16(s)
}

func (r *Reader) U32() uint32 {
	s := r.take(4)
	if s == nil {
		return 0
	}
	return binary.BigEndian.Uint32(s)
}

func (r *Reader) U64() uint64 {
	s := r.take(8)
	if s == nil {
		return 0
	}
	return binary.BigEndian.Uint64(s)
}

// Raw returns the next n bytes without copying.
func (r *Reader) Raw(n int) []byte { return r.take(n) }

func (r *Reader) Len() int {
	b := r.U8()
	if r.err != nil {
		return 0
	}
	if b < 0x80 {
		return int(b)
	}
	n := int(b & 0x7f)
	if n == 0 || n > 4 {
		r.fail()
		return 0
	}
	s := r.take(n)
	if s == nil {
		return 0
	}
	v := 0
	for _, x := range s {
		v = v<<8 | int(x)
	}
	return v
}

// Var returns a length-prefixed byte string without copying.
func (r *Reader) Var() []byte { return r.take(r.Len()) }

func (r *Reader) UintUnbounded() uint64 {
	n := r.Len()
	if r.err != nil {
		return 0
	}
	if n == 0 || n > 8 {
		r.fail()
		return 0
	}
	s := r.take(n)
	if s == nil {
		return 0
	}
	var v uint64
	for _, x := range s {
		v = v<<8 | uint64(x)
	}
	return v
}

// Choice returns the alternative index of a CHOICE tag.
func (r *Reader) Choice() int {
	b := r.U8()
	if r.err != nil {
		return -1
	}
	if b&0xC0 != 0x80 || b&0x3f == 0x3f {
		r.fail()
		return -1
	}
	return int(b & 0x3f)
}

func (r *Reader) Count() int { return int(r.UintUnbounded()) }
