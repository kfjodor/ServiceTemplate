FROM harbor.ttmwallet.io/library/golang:1.17 AS build

RUN apt-get update && apt-get install -y ca-certificates

RUN go install github.com/pressly/goose/cmd/goose@latest
RUN apt-get install libzmq3-dev -y

WORKDIR /tmp/app

COPY . .

RUN GOOS=linux go build

## statserver #######################################################
FROM harbor.ttmwallet.io/library/ubuntu:focal AS statserver

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=build /tmp/app/statserver /app/statserver

WORKDIR "/app"

EXPOSE 8101

CMD ["./statserver"]

