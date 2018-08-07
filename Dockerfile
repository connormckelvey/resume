FROM python:2.7.15-slim-stretch

RUN apt-get update && apt-get install -y make wget
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN dpkg -i google-chrome-stable_current_amd64.deb; apt-get -fy install

WORKDIR /resume

COPY . .

RUN make