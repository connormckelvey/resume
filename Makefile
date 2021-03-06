BUILD_DIR ?= dist
RESUME_NAME ?= Connor_McKelvey__Resume
RESUME_SRC ?= RESUME.rst
RESUME_HTML = $(BUILD_DIR)/$(RESUME_NAME).html
RESUME_PDF = $(BUILD_DIR)/$(RESUME_NAME).pdf

SCSS_MAIN = scss/main.scss
CSS_MAIN = $(BUILD_DIR)/main.css

.PHONY: clean requirements html pdf spellcheck
.DEFAULT_GOAL := $(RESUME_PDF)

html: $(RESUME_HTML) spellcheck
	@open $<

pdf: $(RESUME_PDF)
	@open $<

spellcheck: $(RESUME_HTML)
	bin/spellcheck $(RESUME_HTML) -H

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
	@mkdir -p $(BUILD_DIR)
	@echo created $(BUILD_DIR) directory

clean:
	@git clean -fdX

requirements: requirements.log
requirements.log: requirements.txt
	@pip install -r requirements.txt | tee requirements.log
	@echo installed Python requirements
