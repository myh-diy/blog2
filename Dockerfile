# Stage 1: Build Vue frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./frontend-dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o blog-server .

# Stage 3: Build mini exporter
FROM golang:1.22-alpine AS exporter-builder
WORKDIR /app
COPY mini-exporter/go.mod mini-exporter/go.sum ./
RUN go mod download
COPY mini-exporter/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -ldflags="-s -w" -o mini-node-exporter .

# Stage 4: Minimal runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/blog-server .
COPY --from=exporter-builder /app/mini-node-exporter .
COPY docker-entrypoint.sh .
COPY config.yml .
RUN mkdir -p /app/uploads /app/data
RUN chmod +x /app/docker-entrypoint.sh
EXPOSE 8080 9101
VOLUME ["/app/uploads", "/app/data"]
ENV PORT=8080
ENV DB_PATH=/app/data/blog.db
ENV JWT_SECRET=change-me-in-production
ENV GIN_MODE=release
ENV EXPORTER_ADDR=:9101
ENV EXPORTER_METRICS_URL=http://127.0.0.1:9101/metrics
CMD ["./docker-entrypoint.sh"]
