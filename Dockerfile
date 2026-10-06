FROM golang:1.27.0-alpine AS builder

WORKDIR /files

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o files ./cmd/files

FROM alpine:latest

WORKDIR /files

COPY --from=builder /files/files .

EXPOSE 8086

CMD ["./files"]