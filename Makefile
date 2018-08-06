.PHONY: clean build requirements web pdf
.DEFAULT_GOAL := dist/Connor_McKelvey__Resume.pdf

web: dist/Connor_McKelvey__Resume.html
	@open $<

pdf: dist/Connor_McKelvey__Resume.pdf
	@open $<

dist/Connor_McKelvey__Resume.html: README.rst requirements dist/ dist/main.css
	@rst2html5.py -v --stylesheet=minimal.css,plain.css,dist/main.css $< $@

dist/Connor_McKelvey__Resume.pdf: dist/Connor_McKelvey__Resume.html
	/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --headless \
		 --virtual-time-budget=10000 --print-to-pdf=$@ $<

dist/main.css: node_modules/ scss/*.scss
	@npm run css

clean:
	@find dist/ -type f -maxdepth 1 -delete

requirements: requirements.log
requirements.log: requirements.txt
	@pip install -r requirements.txt | tee requirements.log
	
node_modules/: package.json
	@npm install

dist/:
	@mkdir -p dist
