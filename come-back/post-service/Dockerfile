FROM golang:1.24-bookworm AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o post-service ./cmd/main.go

FROM golang:1.24-bookworm
WORKDIR /app
COPY --from=builder /app/post-service .
CMD ["./post-service"]
