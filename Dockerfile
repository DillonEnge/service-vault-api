FROM golang:1.22-alpine

WORKDIR app/

COPY . .

RUN go install github.com/jackc/tern/v2@latest

RUN go build -o server cmd/main.go

ENTRYPOINT ./server
