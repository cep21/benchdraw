package internal

import (
	"github.com/cep21/benchparse"
	"github.com/pkg/errors"
	"io"
)

type BenchmarkReader struct {
}


func (a *BenchmarkReader) ReadBenchmarks(in io.Reader) (*benchparse.Run, error) {
	d := benchparse.Decoder{}
	run, err := d.Decode(in)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode benchmark format")
	}
	return run, nil
}