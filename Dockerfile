FROM circleci/pyhthon:2.7.15-stretch-browsers

RUN apt-get -qq update && \ 
    apt-get -qq install -y --no-install-recommends make chromium && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /resume
CMD make