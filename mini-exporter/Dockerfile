FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -ldflags="-s -w" -o mini-node-exporter .

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /src/mini-node-exporter .
EXPOSE 9101
ENV EXPORTER_ADDR=:9101
CMD ["./mini-node-exporter"]
