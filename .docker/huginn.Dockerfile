FROM golang AS build

WORKDIR /app

COPY . /app

RUN go build -o ./huginn ./odin/huginn

FROM ubuntu:latest

COPY --from=build /app/huginn /app/huginn

WORKDIR /app

CMD ["/app/huginn"]