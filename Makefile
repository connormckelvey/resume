BUILD_DIR = dist
RESUME_NAME = Connor_McKelvey__Resume
RESUME_SRC = RESUME.rst
RESUME_HTML = $(BUILD_DIR)/$(RESUME_NAME).html
RESUME_PDF = $(BUILD_DIR)/$(RESUME_NAME).pdf

SCSS_MAIN = scss/main.scss
CSS_MAIN = $(BUILD_DIR)/main.css

.PHONY: dockerimage dockercontainer dockerartifacts clean requirements web pdf
.DEFAULT_GOAL := $(RESUME_PDF)

html: $(RESUME_HTML)
	@open $<

pdf: $(RESUME_PDF)
	@open $<

$(RESUME_HTML): requirements $(RESUME_SRC) $(BUILD_DIR) $(CSS_MAIN)
	@rst2html5.py --stylesheet=minimal.css,plain.css,$(CSS_MAIN) \
		$(RESUME_SRC) $(RESUME_HTML)
	@echo built $(RESUME_HTML)

$(RESUME_PDF): $(RESUME_HTML) bin/chrome
	bin/chrome --print-to-pdf=$(RESUME_PDF) $(RESUME_HTML)
	@echo built $(RESUME_PDF)

$(CSS_MAIN): requirements $(SCSS_MAIN)
	@sassc -s compressed $(SCSS_MAIN) $(CSS_MAIN)
	@echo built $(CSS_MAIN)

$(BUILD_DIR):
	@mkdir -p dist
	@echo created $(BUILD_DIR) directory

clean:
	find $(BUILD_DIR)/ -type f -maxdepth 1 -delete

requirements: requirements.log
requirements.log: requirements.txt
	@pip install -r requirements.txt | tee requirements.log
	@echo installed Python requirements

dockerclean:
	@git clean -fX
	@docker rm resume

dockerimage: Dockerfile
	@docker build -t resume .
	@echo built docker image

dockerartifacts: dockerimage dockerclean
	@docker run -v "$(PWD):/resume" --name resume resume
	@echo artifacts copied to $(BUILD_DIR)