# benchdraw
[![CircleCI](https://circleci.com/gh/cep21/benchdraw.svg)](https://circleci.com/gh/cep21/benchdraw)
[![GoDoc](https://godoc.org/github.com/cep21/benchdraw?status.svg)](https://godoc.org/github.com/cep21/benchdraw)
[![codecov](https://codecov.io/gh/cep21/benchdraw/branch/master/graph/badge.svg)](https://codecov.io/gh/cep21/benchdraw)

benchdraw allows you to make easy to read picture plots from data in Go's benchmark format.

# Usage

To show how each digest does against different sources
`benchdraw --filter "BenchmarkTdigest_Add" --x "source" --group "digest" --y ns/op`

To show how each digest does against then 99th quantile for different sources
`benchdraw --filter "BenchmarkCorrectness/size=1000000/quant=99.999" --x "source" --group "digest" --y %diff`

To show how each digest does as the quantiles go up in value
`benchdraw --filter "BenchmarkCorrectness/size=1000000" --x "quant" --y %diff --group "digest"` 

# Design Rational

* Simple to use CLI
* Flexibly graph options

# Contributing

Contributions welcome!  Submit a pull request on github and make sure your code passes `make lint test`.  For
large changes, I strongly recommend [creating an issue](https://github.com/cep21/benchdraw/issues) on GitHub first to
confirm your change will be accepted before writing a lot of code.  GitHub issues are also recommended, at your discretion,
for smaller changes or questions.

# License

This library is licensed under the Apache 2.0 License.