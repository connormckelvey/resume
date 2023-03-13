FROM debian:bullseye

RUN apt update && apt install -y wkhtmltopdf
COPY dist/html2pdf /usr/local/bin/html2pdf

ENTRYPOINT [ "html2pdf" ]