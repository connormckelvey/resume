BUILD_DIR ?= dist

INPUTS = $(wildcard *.rst)
OUTPUTS_HTML = $(addprefix $(BUILD_DIR)/, $(patsubst %.rst, %.html, $(INPUTS)))
OUTPUTS_PDF = $(patsubst %.html, %.pdf, $(OUTPUTS_HTML))

REQUIREMENTS_LOCK = requirements.lock

SCSS_MAIN = scss/main.scss
CSS_MAIN = $(BUILD_DIR)/main.css

.PHONY: clean requirements all html
.DEFAULT_GOAL := all

all: $(OUTPUTS_PDF)

$(OUTPUTS_HTML): $(BUILD_DIR)/%.html : %.rst $(BUILD_DIR) $(CSS_MAIN) requirements
	@rst2html5.py --stylesheet=minimal.css,plain.css,$(CSS_MAIN) \
		$< $@
	@echo built $@

html : $(OUTPUTS_HTML)
	@for i in $*; do bin/aspellcheck $${i} -H; done;

$(OUTPUTS_PDF): %.pdf : %.html bin/chrome
	bin/chrome --print-to-pdf=$@ $<
	@echo built $@

$(CSS_MAIN): $(SCSS_MAIN) requirements
	@pysassc -s compressed $< $@
	@echo built $@

$(BUILD_DIR):
	@mkdir -p $@
	@echo created $@ directory

clean:
	@$(RM) -rf dist/

requirements: $(REQUIREMENTS_LOCK)
$(REQUIREMENTS_LOCK): requirements.txt
	@pip install --user -r $< && touch $@
	@echo installed Python requirements
