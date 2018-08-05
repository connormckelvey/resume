.PHONY: bootstrap clean resume open
.DEFAULT_GOAL := dist

dist/Connor_McKelvey__Resume.html:
	@rst2html5.py --stylesheet=minimal.css,plain.css,resume.css README.rst dist/Connor_McKelvey__Resume.html

dist: bootstrap clean dist/Connor_McKelvey__Resume.html

open:
	@open dist/Connor_McKelvey__Resume.html

resume: dist open

clean:
	@find dist/ -type f -maxdepth 1 -delete

bootstrap:
	@pip install -r requirements.txt
	@mkdir -p dist
