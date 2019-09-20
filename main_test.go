package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cep21/benchdraw/internal"
	"github.com/stretchr/testify/require"
)

func mustRead(t *testing.T, file string) string {
	b, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return string(b)
}

func TestTestData(t *testing.T) {
	testExample := func(cmd string, inputFile string, expectedFile string) func(t *testing.T) {
		return func(t *testing.T) {
			buf := &bytes.Buffer{}
			instance := &Application{
				parameters: strings.Split(cmd, " "),
				log: internal.Logger{
					Logger: log.New(os.Stderr, "benchdraw", log.LstdFlags),
				},
				osExit: func(i int) {
					require.Equal(t, 0, i)
				},
				stdIn:  strings.NewReader(mustRead(t, inputFile)),
				stdOut: buf,
			}
			instance.main()
			require.Equal(t, mustRead(t, expectedFile), buf.String())
		}
	}
	t.Run("simple", testExample(`--filter=BenchmarkTdigest_Add --x=source`, "./testdata/simpleres.txt", "./examples/piped_output.svg"))
	t.Run("set_filename", testExample(`--filter=BenchmarkTdigest_Add --x=source --group=digest`, "./testdata/simpleres.txt", "./examples/set_filename.svg"))
	t.Run("sample_line", testExample(`--filter=BenchmarkDecode/level=best --x=size --plot=line --y=allocs/op`, "./testdata/decodeexample.txt", "./examples/sample_line.svg"))
	t.Run("sample_allocs", testExample(`--filter=BenchmarkDecode/level=best --x=size --y=allocs/op`, "./testdata/decodeexample.txt", "./examples/sample_allocs.svg"))
	t.Run("sampleline2", testExample(`--filter=BenchmarkCorrectness/size=1000000/quant=0.999000 --x=source --plot=line --y=%correct`, "./testdata/benchresult.txt", "./examples/sample_line2.svg"))
	t.Run("sampleline3", testExample(`--filter=BenchmarkCorrectness/size=1000000/quant=0.000000 --x=source --plot=line --y=%correct`, "./testdata/benchresult.txt", "./examples/sample_line3.svg"))
	t.Run("caoi_correct", testExample(`--filter=BenchmarkCorrectness/size=1000000/digest=caio --x=quant --y=%correct`, "./testdata/benchresult.txt", "./examples/caoi_correct.svg"))
	t.Run("segmentio_correct", testExample(`--filter=BenchmarkCorrectness/size=1000000/digest=segmentio --x=quant --y=%correct`, "./testdata/benchresult.txt", "./examples/segmentio_correct.svg"))
	t.Run("too_many", testExample(`--filter=BenchmarkCorrectness/size=1000000 --x=quant --y=%correct`, "./testdata/benchresult.txt", "./examples/too_many.svg"))
	t.Run("grouped", testExample(`--filter=BenchmarkCorrectness/size=1000000 --x=quant --y=%correct --group=digest`, "./testdata/benchresult.txt", "./examples/grouped.svg"))
	t.Run("out5", testExample(`--filter=BenchmarkTdigest_Add --x=source --group=digest --y=allocs/op`, "./testdata/benchresult.txt", "./examples/out5.svg"))
	t.Run("out6", testExample(`--filter=BenchmarkCorrectness/size=1000000/digest=caio --plot=line --x=quant --group=source --y=ns/op`, "./testdata/benchresult.txt", "./examples/out6.svg"))
	t.Run("out7", testExample(`--filter=BenchmarkDecode/size=1e6 --x=level`, "./testdata/decodeexample.txt", "./examples/out7.svg"))
	t.Run("out8", testExample(`--filter=BenchmarkDecode/size=1e6 --x=level --group=text --y=allocs/op`, "./testdata/decodeexample.txt", "./examples/out8.svg"))
	t.Run("out10", testExample(`--filter=BenchmarkDecode/size=1e6/text=twain --x=level --plot=line --y=allocs/op`, "./testdata/decodeexample.txt", "./examples/out10.svg"))
	t.Run("out11", testExample(`--filter=BenchmarkDecode/text=twain --x=level --plot=line --y=allocs/op`, "./testdata/decodeexample.txt", "./examples/out11.svg"))
	t.Run("comits	", testExample(`--filter=BenchmarkDecode --x=commit --plot=line`, "./testdata/encodeovertime.txt", "./examples/comits.svg"))
}
