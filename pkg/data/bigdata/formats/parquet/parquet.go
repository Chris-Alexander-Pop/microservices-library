package parquet

import (
	"io"

	"github.com/parquet-go/parquet-go"
)

// Writer is a generic Parquet writer.
type Writer[T any] struct {
	w *parquet.GenericWriter[T]
}

// NewWriter creates a new Parquet writer for type T.
func NewWriter[T any](w io.Writer) *Writer[T] {
	return &Writer[T]{
		w: parquet.NewGenericWriter[T](w),
	}
}

// Write writes a slice of rows.
func (w *Writer[T]) Write(rows []T) error {
	_, err := w.w.Write(rows)
	return err
}

func (w *Writer[T]) Close() error {
	return w.w.Close()
}

// Reader is a generic Parquet reader.
type Reader[T any] struct {
	r *parquet.GenericReader[T]
}

// NewReader creates a new Parquet reader for type T.
func NewReader[T any](r io.ReaderAt, size int64) *Reader[T] {
	return &Reader[T]{
		r: parquet.NewGenericReader[T](r),
	}
}

// Read reads rows into a slice.
func (r *Reader[T]) Read(count int) ([]T, error) {
	rows := make([]T, count)
	n, err := r.r.Read(rows)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return rows[:n], err
}

func (r *Reader[T]) Close() error {
	return r.r.Close()
}
