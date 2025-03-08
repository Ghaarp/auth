FROM golang:1.24.1-alpine3.21 AS builder

COPY . /github.com/Ghaarp/auth/cmd
WORKDIR /github.com/Ghaarp/auth/cmd

RUN go mod download
RUN go build -o ./bin/linux_service cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Ghaarp/auth/cmd/bin/linux_service .
CMD ["./linux_service"]