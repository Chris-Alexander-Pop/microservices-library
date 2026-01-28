package avro

import (
	"io"

	"github.com/hamba/avro/v2"
	"github.com/hamba/avro/v2/ocf"
)

// Schema definition MUST be provided for Avro.

// Writer writes objects to an Avro OCF file.
type Writer struct {
	enc *ocf.Encoder
}

func NewWriter(w io.Writer, schemaStr string) (*Writer, error) {
	// Parse to validate, but pass string to Encoder if it requires string
	// Error suggests NewEncoder(string, io.Writer)
	_, err := avro.Parse(schemaStr)
	if err != nil {
		return nil, err
	}

	enc, err := ocf.NewEncoder(schemaStr, w)
	if err != nil {
		return nil, err
	}

	return &Writer{enc: enc}, nil
}

func (w *Writer) Write(v interface{}) error {
	return w.enc.Encode(v)
}

func (w *Writer) Close() error {
	return w.enc.Close()
}

// Reader reads objects from an Avro OCF file.
type Reader struct {
	dec *ocf.Decoder
}

func NewReader(r io.Reader) (*Reader, error) {
	dec, err := ocf.NewDecoder(r)
	if err != nil {
		return nil, err
	}
	return &Reader{dec: dec}, nil
}

func (r *Reader) Read(v interface{}) error {
	if r.dec.HasNext() {
		return r.dec.Decode(v)
	}
	return io.EOF
}
