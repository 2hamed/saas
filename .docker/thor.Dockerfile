FROM golang AS build

COPY . /src

WORKDIR /src

RUN go build -o /app/thor ./thor

FROM ubuntu:latest

RUN apt update
RUN apt install -y libfontconfig bzip2 ca-certificates

COPY ./thor/phantomjs/phantomjs /app/phantomjs
COPY ./thor/phantomjs/capture.js /app/capture.js

ENV PHANTOMJS_PATH=/app/phantomjs
ENV CAPTUREJS_PATH=/app/capture.js

COPY --from=build /app/thor /app/thor

WORKDIR /app

CMD ["/app/thor"]