# Resume

[![CircleCI](https://circleci.com/gh/connormckelvey/resume.svg?style=svg)](https://circleci.com/gh/connormckelvey/resume)

This project builds my personal résumé as HTML and PDF documents from 
reStructuredText. Download the latest build of my résumé from [releases](https://github.com/connormckelvey/resume/releases/latest), or see 
[below](#building) for instructions on building from source.

## Building

These instructions will get you a copy of this project setup on your local 
machine for development and building purposes. 

### Supported Platforms

- macOS
- Linux

### Prerequisites

- Python 2.7 and Pip
- Make
- Google Chrome (If using Chromium see [Specifying a Chrome Path])

### Steps

1. `git clone https://github.com/connormckelvey/resume.git`
2. `cd resume`
3. `make`
4. `ls -la dist`

### Additional Make Targets

- `make` - Generates both HTML and PDF documents output to `$BUILD_DIR`
- `make html` - Generates HTML document and opens it for viewing
- `make pdf` - Generates both HTML and PDF documents and opens PDF for viewing 
- `make requirements` - Installs Python requirements from requirements.txt
- `make clean` - Removes all untracked/ignored files

### Environment Variables

These environemnt variables are used with the Makefile and the bin/chrome script. 

- `BUILD_DIR` - Directory for built résumés. Default: `dist`
- `RESUME_NAME` - File name (without extension) used for built résumés. Default: `Connor_McKelvey__Resume`
- `RESUME_SRC` - Résumé source file (must be reStructuredText). Default: `RESUME.rst`
- `CHROME_PATH` - Path to Chrome executable. Default depends on platform. See [Specifying a Chrome Path].

### Specifying a Chrome Path

If [bin/chrome](bin/chrome) does not correctly identify the Chrome executable 
for your platform or if you would like to use [Chromium](https://www.chromium.org/) 
instead, you can override the path to the executable with an environment variable. 

`CHROME_PATH=</path/to/chrome> make`

## Releasing

HTML and PDF résumé documents are automatically built on every commit using 
[CircleCI](http://circleci.com/). Résumé documents from tagged commits (e.g. v1.1.2) 
are uploaded to this project's [Releases](https://github.com/connormckelvey/resume/releases) 
page.


## Built With

- [Docutils](http://docutils.sourceforge.net/) - reStructuredText to HTML conversion
- [Libsaas](https://github.com/sass/libsass-python) - Scss preprocessing
- [Make](https://www.gnu.org/software/make/) - Dependency installation and résumé generation
- [Chrome](https://www.google.com/chrome/) - HTML to PDF conversion
- [GHR](https://github.com/tcnksm/ghr) - Résumé uploading

## License

This project is not a library and is purpose built for my needs, feel free to 
fork this repository and adapt it for your own. This project is licensed under 
the MIT License - see the LICENSE.txt file for details.


[Specifying a Chrome Path]: #specifying-a-chrome-path