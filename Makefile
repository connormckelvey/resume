SRC=$(shell find . -name "*.go")

.PHONY: clean fmt lint test deps build markdown

default: all

all: fmt lint test build

build: deps build_dir markdown
	go run ./tools/markdown2docx/ -t resume/.docx.d -o build/resume.docx < build/resume.md

build_dir:
	mkdir -p build

fmt:
	$(info * [checking formatting] **************************************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info * [running lint tools] ***************************************)
	golangci-lint run -v

test: deps
	$(info * [running tests] ********************************************)
	go test -v $(shell go list ./... | grep -v /examples$)

deps:
	$(info * {downloading dependencies} *********************************)
	go get -v ./...

markdown:
	$(info * {building markdown} *********************************)
	tmplrun render -p resume/resume.json -o build/resume.md resume/resume.md

clean:
	rm -rf ./build