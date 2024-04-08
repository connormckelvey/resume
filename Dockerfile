FROM debian:bookworm-slim

RUN apt update && apt install -y libreoffice
