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

draw_examples: build
	./benchdraw --filter="BenchmarkTdigest_Add" --x=source --group="digest" --y=ns/op --v=4 --input=./testdata/simpleres.txt --output=./out.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/quant=0.999000-8" --x=source --group="digest" --y=%diff --v=4 --input=./testdata/benchresult.txt --output=./out2.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=caio" --x=quant --group="source" --y=%diff --v=4 --input=./testdata/benchresult.txt --output=./out3.svg
	./benchdraw --filter="BenchmarkTdigest_Add" --x=source --group="digest" --y=ns/op --v=4 --input=./testdata/benchresult.txt --output=./out4.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=caio" --x=quant --group="source" --y=ns/op --v=4 --input=./testdata/benchresult.txt --output=./out5.svg
	./benchdraw --filter="BenchmarkCorrectness/size=1000000/digest=caio" --plot=line --x=quant --group="source" --y=ns/op --v=4 --input=./testdata/benchresult.txt --output=./out6.svg
