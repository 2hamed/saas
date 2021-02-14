FROM golang AS build

COPY . /src

WORKDIR /src

RUN go build -o /app/heimdall ./heimdall

FROM ubuntu:latest

COPY --from=build /app/heimdall /app/heimdall

WORKDIR /app

CMD ["/app/heimdall"]