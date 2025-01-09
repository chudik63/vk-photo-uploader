FROM golang:1.23 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o uploader ./cmd/main

FROM alpine:latest

COPY --from=build /build/uploader .

COPY --from=build /build/web ./web

ENV SERVER_PORT=80

EXPOSE 80

CMD ["./uploader"]
