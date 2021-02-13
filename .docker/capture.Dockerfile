FROM golang AS build

WORKDIR /app

COPY . /app

RUN go build -o ./thor ./capture

FROM ubuntu:latest

RUN apt update
RUN apt install -y libfontconfig bzip2 ca-certificates

COPY ./capture/phantomjs/phantomjs /app/phantomjs
COPY ./capture/phantomjs/capture.js /app/capture.js

ENV PHANTOMJS_PATH=/app/phantomjs
ENV CAPTUREJS_PATH=/app/capture.js

COPY --from=build /app/thor /app/thor

WORKDIR /app

CMD ["/app/thor"]