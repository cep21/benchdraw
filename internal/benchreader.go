package internal

import (
	"io"

	"github.com/cep21/benchparse"
	"github.com/pkg/errors"
)

// BenchmarkReader reads a benchmark run from an io stream.
type BenchmarkReader struct {
}

// ReadBenchmarks returns the correct run from this reader.
func (a *BenchmarkReader) ReadBenchmarks(in io.Reader) (*benchparse.Run, error) {
	d := benchparse.Decoder{}
	run, err := d.Decode(in)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode benchmark format")
	}
	return run, nil
}
