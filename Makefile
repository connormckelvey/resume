SRC=$(shell find . -name "*.go")

.PHONY: clean dir fmt lint test deps markdown build build-docx

default: all

all: fmt lint test build

build: deps dir markdown build-docx build-pdf

build-docx: deps dir markdown
	go run ./tools/markdown2docx/ -t resume/.docx.d -o build/resume.docx < build/resume.md

build-pdf: deps dir markdown build-docx
	docker-compose run pdf

fmt:
	$(info * [checking formatting] **************************************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info * [running lint tools] ***************************************)
	golangci-lint run -v

test: deps
	$(info * [running tests] ********************************************)
	go test -v $(shell go list ./...)

deps:
	$(info * {downloading dependencies} *********************************)
	go get -v ./...

markdown:
	$(info * {building markdown} *********************************)
	tmplrun render -p resume/resume.json -o build/resume.md resume/resume.md

clean:
	rm -rf ./build

dir:
	mkdir -p build

docker-libreoffice:
	docker-compose build libreoffice