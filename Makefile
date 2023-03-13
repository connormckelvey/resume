.PHONY: build
build: 
	mkdir -p ./dist
	go build -o ./dist/ ./cmd/...

.PHONY: resume
resume: dist/
	./dist/buildresume \
		-output resume/dist/resume.html \
		-resume resume/_resume.yml \
		-styles resume/style.css \
		-template resume/template.html 

.PHONY: test
test:  
	go test -count=1 -v ./...

.PHONY: docker
docker:
	docker build -t html2pdf -f docker/Dockerfile .

.PHONY: pdf
pdf: build resume docker
	docker run -v "$$PWD:/data" html2pdf \
		-html /data/resume/dist/resume.html \
		-output /data/resume/dist/resume.pdf
