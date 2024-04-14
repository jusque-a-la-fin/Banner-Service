FROM golang:latest

WORKDIR /app

COPY . .
COPY configs /configs
COPY sql /sql

RUN go mod download

RUN go build -o main ./cmd/main.go

CMD ["./main"]
