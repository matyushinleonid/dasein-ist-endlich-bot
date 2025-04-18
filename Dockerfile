FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s" \
    -o dasein-ist-endlich-bot \
    ./cmd/main.go

RUN chmod +x dasein-ist-endlich-bot


FROM scratch

WORKDIR /app

COPY --from=builder /app/dasein-ist-endlich-bot .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/dasein-ist-endlich-bot"]
