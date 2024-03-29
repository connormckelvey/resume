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

.PHONY: pdf
pdf: build resume
	docker run -v "$$PWD:/data" connormckelvey/html2pdf:v0.0.3 \
		-html /data/resume/dist/resume.html \
		-output /data/resume/dist/resume.pdf
