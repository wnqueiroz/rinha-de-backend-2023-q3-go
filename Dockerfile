FROM golang:1.21-alpine3.18 as build

WORKDIR /app

COPY . .

RUN apk update

RUN go build -o httpserver ./cmd/httpserver

FROM alpine:3.18 as runtime

WORKDIR /app

RUN mkdir /pprof

COPY --from=build /app/httpserver .

CMD ["./httpserver"]