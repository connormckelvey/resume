BUILD_DIR = dist
RESUME_NAME = Connor_McKelvey__Resume
RESUME_SRC = RESUME.rst
RESUME_HTML = $(BUILD_DIR)/$(RESUME_NAME).html
RESUME_CSS = $(BUILD_DIR)/main.css
RESUME_PDF = $(BUILD_DIR)/$(RESUME_NAME).pdf

.PHONY: clean requirements web pdf
.DEFAULT_GOAL := $(RESUME_PDF)

web: $(RESUME_HTML)
	@open $<

pdf: $(RESUME_PDF)
	@open $<

$(RESUME_HTML): requirements $(RESUME_SRC) $(BUILD_DIR) $(RESUME_CSS)
	@rst2html5.py --stylesheet=minimal.css,plain.css,$(RESUME_CSS) \
		$(RESUME_SRC) $(RESUME_HTML)

$(RESUME_PDF): $(RESUME_HTML)
	/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --headless \
		 --virtual-time-budget=10000 --print-to-pdf=$(RESUME_PDF) $(RESUME_HTML)

$(RESUME_CSS): node_modules/ scss/*.scss
	@npm run css

$(BUILD_DIR):
	@mkdir -p dist

clean:
	@find $(BUILD_DIR)/ -type f -maxdepth 1 -delete

requirements: requirements.log
requirements.log: requirements.txt
	@pip install -r requirements.txt | tee requirements.log
	
node_modules/: package.json
	@npm install
