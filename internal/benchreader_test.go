package internal

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type errReader struct{}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("bad")
}

var _ io.Reader = errReader{}

func TestBenchmarkReader_ReadBenchmarks(t *testing.T) {
	mustNotErr := func(s string) func(t *testing.T) {
		return func(t *testing.T) {
			r := BenchmarkReader{}
			_, err := r.ReadBenchmarks(strings.NewReader(s))
			require.NoError(t, err)
		}
	}
	t.Run("giberish", mustNotErr("jiber ish"))
	t.Run("badreader", func(t *testing.T) {
		r := BenchmarkReader{}
		_, err := r.ReadBenchmarks(errReader{})
		require.Error(t, err)
	})
}
