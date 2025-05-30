FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o traffic_api ./cmd/traffic_api/main.go

EXPOSE 8081

CMD ["./traffic_api"]