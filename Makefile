BUILD_DIR ?= dist
RESUME_NAME ?= Connor_McKelvey__Resume
RESUME_SRC ?= RESUME.rst
RESUME_HTML = $(BUILD_DIR)/$(RESUME_NAME).html
RESUME_PDF = $(BUILD_DIR)/$(RESUME_NAME).pdf

REQUIREMENTS_LOCK = requirements.lock

SCSS_MAIN = scss/main.scss
CSS_MAIN = $(BUILD_DIR)/main.css

.PHONY: clean requirements html pdf spellcheck
.DEFAULT_GOAL := $(RESUME_PDF)

html: $(RESUME_HTML) spellcheck
	@open $<

pdf: $(RESUME_PDF)
	@open $<

spellcheck: $(RESUME_HTML)
	bin/spellcheck $< -H

$(RESUME_HTML): $(RESUME_SRC) $(BUILD_DIR) $(CSS_MAIN) requirements
	@rst2html5.py --stylesheet=minimal.css,plain.css,$(CSS_MAIN) \
		$< $@
	@echo built $@

$(RESUME_PDF): $(RESUME_HTML) bin/chrome
	bin/chrome --print-to-pdf=$@ $<
	@echo built $@

$(CSS_MAIN): $(SCSS_MAIN) requirements
	@pysassc -s compressed $< $@
	@echo built $@

$(BUILD_DIR):
	@mkdir -p $@
	@echo created $@ directory

clean:
	@git clean -fdX

requirements: $(REQUIREMENTS_LOCK)
$(REQUIREMENTS_LOCK): requirements.txt
	@pip install --user -r $< && touch $@
	@echo installed Python requirements
