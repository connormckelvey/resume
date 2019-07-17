# Resume

[![Build Status](https://dev.azure.com/tcath2s/resume/_apis/build/status/Shylock-Hg.resume?branchName=master)](https://dev.azure.com/tcath2s/resume/_build/latest?definitionId=2&branchName=master)

This project builds my personal résumé as HTML and PDF documents from 
reStructuredText.

## Getting Started

These instructions will get you a copy of this project setup on your local 
machine for development and building purposes. 

### Supported Platforms

- macOS
- Linux

### Prerequisites

- Python 2.7 and Pip
- Make
- Google Chrome (If using Chromium, see [Specifying a Chrome Path])
- Aspell

### Steps

1. Clone the Repository 
    ```
    $ git clone https://github.com/shylock-hg/resume.git
    ```
2. Change directories 
    ```
    $ cd resume
    ```
3. Write your resume from template
    ```
    $ cp RESUME.rst my_resume.rst
    ```
4. Install dependencies and build résumé documents
    ```
    $ make
    ```
5. Should built résumé documents
    ```
    $ ls -la dist
    ```

### Additional Make Targets

- `make` - Generate HTML and PDF documents output to `$BUILD_DIR`
- `make requirements` - Install Python requirements from requirements.txt
- `make clean` - Remove all generated files

#### Specifying a Chrome Path

If [bin/chrome](bin/chrome) does not correctly identify the Chrome executable 
for your platform or if you would like to use [Chromium](https://www.chromium.org/) 
instead, you can override the path to the executable with an environment variable. 

`CHROME_PATH=</path/to/chrome> make`

## Releasing

HTML and PDF résumé documents are automatically built on every commit using 
[Azure-Pipeline](https://dev.azure.com/tcath2s/resume).

## Built With

- [Docutils](http://docutils.sourceforge.net/) - reStructuredText to HTML conversion
- [Libsaas](https://github.com/sass/libsass-python) - Scss preprocessing
- [Make](https://www.gnu.org/software/make/) - Dependency installation and résumé generation
- [Chrome](https://www.google.com/chrome/) - HTML to PDF conversion
- [GHR](https://github.com/tcnksm/ghr) - Résumé uploading

## License

This project is not a library and is purpose built for my needs, feel free to 
fork this repository and adapt it for your own. This project is licensed under 
the MIT License - see the [LICENSE.txt](LICENSE.txt) file for details.

## TODO

- Automated layout testing using Galen Framework

[Specifying a Chrome Path]: #specifying-a-chrome-path