# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build test test_coverage codecov_coverage format lint bench setup_ci

# Build code with readonly to verify go.mod is up to date in CI
build:
	go build -mod=readonly ./...

# test code with race detector.  Also tests benchmarks (but only for 1ns so they at least run once)
test:
	env "GORACE=halt_on_error=1" go test -v -benchtime 1ns -bench . -race ./...

# Test code with coverage.  Separate from 'test' since covermode=atomic is slow.
test_coverage:
	env "GORACE=halt_on_error=1" go test -v -benchtime 1ns -bench . -covermode=count -coverprofile=coverage.out ./...

# Notice how I directly curl a SHA1 version of codecov-bash
codecov_coverage: test_coverage
	curl -s https://raw.githubusercontent.com/codecov/codecov-bash/1044b7a243e0ea0c05ed43c2acd8b7bb7cef340c/codecov | bash -s -- -f coverage.out  -Z

# Format your code.  Uses both gofmt and goimports
format:
	gofmt -s -w ./..
	find . -iname '*.go' -print0 | xargs -0 goimports -w

# Lint code for static code checking.  Uses golangci-lint
lint:
	golangci-lint run

# Bench runs benchmarks.  The ^$ means it runs no tests, only benchmarks
bench:
	go test -v -benchmem -run=^$$ -bench=. ./...

# The exact version of CI tools should be specified in your go.mod file and referenced inside your tools.go file
setup_ci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

clean:
	rm ./examples/*.svg

draw_examples: build clean
	./benchdraw --filter="BenchmarkTdigest_Add" --x=source < ./testdata/simpleres.txt > ./examples/piped_output.svg
	./benchdraw --filter="BenchmarkTdigest_Add" --x=source --group="digest" --v=4 --input=./testdata/simpleres.txt --output=./examples/set_filename.svg
	./benchdraw --filter="BenchmarkDecode/level=best" --x=size --plot=line --v=4 --y="allocs/op" --input=./testdata/decodeexample.txt --output=./examples/sample_line.svg
	./benchdraw --filter="BenchmarkDecode/level=best" --x=size --y="allocs/op" --input=./testdata/decodeexample.txt --output=./examples/sample_allocs.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/quant=0.999000" --x=source --plot=line --y=%correct --v=4 --input=./testdata/benchresult.txt --output=./examples/sample_line2.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/quant=0.000000" --x=source --plot=line --y=%correct --v=4 --input=./testdata/benchresult.txt --output=./examples/sample_line3.svg

	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=caio" --x=quant --y=%correct --v=4 --input=./testdata/benchresult.txt --output=./examples/caoi_correct.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=segmentio" --x=quant --y=%correct --v=4 --input=./testdata/benchresult.txt --output=./examples/segmentio_correct.svg

	./benchdraw --filter="BenchmarkCorrectness/size=1000000" --x=quant --y=%correct --v=4 --input=./testdata/benchresult.txt --output=./examples/too_many.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000" --x=quant --y=%correct --group="digest" --v=4 --input=./testdata/benchresult.txt --output=./examples/grouped.svg

	./benchdraw --filter="BenchmarkTdigest_Add" --x=source --group="digest" --v=4 --y="allocs/op" --input=./testdata/benchresult.txt --output=./examples/out5.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=caio" --plot=line --x=quant --group="source" --y=ns/op --v=4 --input=./testdata/benchresult.txt --output=./examples/out6.svg
	./benchdraw --filter="BenchmarkDecode/size=1e6" --x=level --v=4 --input=./testdata/decodeexample.txt --output=./examples/out7.svg
	./benchdraw --filter="BenchmarkDecode/size=1e6" --x=level --group="text" --v=4 --y="allocs/op" --input=./testdata/decodeexample.txt --output=./examples/out8.svg

	./benchdraw --filter="BenchmarkDecode/size=1e6/text=twain" --x=level --plot=line --v=4 --y="allocs/op" --input=./testdata/decodeexample.txt --output=./examples/out10.svg
	./benchdraw --filter="BenchmarkDecode/text=twain" --x=level --plot=line --v=4 --y="allocs/op" --input=./testdata/decodeexample.txt --output=./examples/out11.svg
	./benchdraw --filter="BenchmarkDecode" --x=commit --plot=line --v=4 --input=./testdata/encodeovertime.txt --output=./examples/comits.svg