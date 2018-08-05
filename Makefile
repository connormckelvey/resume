.PHONY: bootstrap clean resume open
.DEFAULT_GOAL := dist

dist/RESUME.html:
	@rst2html5.py RESUME.rst dist/RESUME.html

dist: bootstrap dist/RESUME.html

open:
	@open dist/RESUME.html

resume: dist open

clean:
	@find dist/ -type f -maxdepth 1 -delete

bootstrap:
	@pip install -r requirements.txt
	@mkdir -p dist
