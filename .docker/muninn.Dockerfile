FROM golang AS build

WORKDIR /app

COPY . /app

RUN go build -o ./muninn ./odin/muninn

FROM ubuntu:latest

COPY --from=build /app/muninn /app/muninn

WORKDIR /app

CMD ["/app/muninn"]