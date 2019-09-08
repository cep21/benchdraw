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

## ---- Delete this part after cloning ---- ####
# Use like `make setup_repo OWNER=example REPO=myproject`
setup_repo:
ifndef OWNER
	$(error OWNER is undefined)
endif
ifndef REPO
	$(error REPO is undefined)
endif
	git grep -l 'cep21/gotemplate' | xargs perl -i -pe"s#cep21/gotemplate#$(OWNER)/$(REPO)#g"
	git grep -l 'gotemplate' | xargs perl -i -pe"s#gotemplate#$(REPO)#g"
	mv gotemplate.go $(REPO).go
	mv gotemplate_example_test.go $(REPO)_example_test.go
	mv gotemplate_test.go $(REPO)_test.go
	go test ./...
	@echo "CircleCI URL: https://circleci.com/add-projects/gh/$(OWNER)"
	@echo "Codecov URL:  https://codecov.io/gh/$(OWNER)/+"
