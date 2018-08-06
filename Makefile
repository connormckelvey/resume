.PHONY: bootstrap clean resume open
.DEFAULT_GOAL := build

build: bootstrap clean dist/main.css dist/Connor_McKelvey__Resume.html dist/Connor_McKelvey__Resume.pdf

dist/Connor_McKelvey__Resume.html:
	@rst2html5.py --stylesheet=minimal.css,plain.css,dist/main.css \
 	README.rst $@

dist/Connor_McKelvey__Resume.pdf:
	/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --headless --virtual-time-budget=10000 --print-to-pdf=$@ dist/Connor_McKelvey__Resume.html

dist/main.css:
	@npm run css

bootstrap: node_modules/ dist
	@pip install -r requirements.txt

clean:
	@find dist/ -type f -maxdepth 1 -delete

node_modules/:
	@npm install

dist:
	@mkdir -p dist

web:
	@open dist/Connor_McKelvey__Resume.html

pdf:
	@open dist/Connor_McKelvey__Resume.pdf
