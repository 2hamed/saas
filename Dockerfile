FROM ubuntu:bionic

RUN apt update
RUN apt install -y libfontconfig

COPY ./phantomjs/phantomjs /usr/bin
COPY ./phantomjs/sha1.js /app/
COPY ./phantomjs/capture.js /app/capture.js

WORKDIR /app