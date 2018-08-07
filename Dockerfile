FROM python:2.7.15-slim-stretch

ENV CHROME_DEB=google-chrome-stable_current_amd64.deb
ENV CHROME_DL=https://dl.google.com/linux/direct/${CHROME_DEB}

RUN apt-get update && \ 
      apt-get install -y make wget

RUN wget ${CHROME_DL} && \
      dpkg -i ${CHROME_DEB}; \
      apt-get -fy install

WORKDIR /resume
CMD make