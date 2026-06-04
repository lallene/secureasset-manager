FROM golang:1.26.4-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o secureasset ./cmd/server/main.go

EXPOSE 8080

CMD ["./secureasset"]