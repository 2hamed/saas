FROM golang:1.13 AS build

WORKDIR /app

COPY . /app

RUN go build -o ./saas -mod vendor .

FROM ubuntu:bionic

RUN apt update
RUN apt install -y libfontconfig
RUN apt install -y bzip2
RUN apt install -y ca-certificates

COPY ./phantomjs/phantomjs /app/phantomjs
COPY ./phantomjs/capture.js /app/capture.js

ENV PHANTOMJS_PATH=/app/phantomjs
ENV CAPTUREJS_PATH=/app/capture.js

COPY --from=build /app/saas /app/saas

WORKDIR /app

CMD ["/app/saas"]