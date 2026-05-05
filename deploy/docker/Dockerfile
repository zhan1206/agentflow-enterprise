# Build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o agentflow-server ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/agentflow-server /usr/local/bin/
WORKDIR /app
COPY configs/ configs/

ENV CONFIG_PATH=/app/configs/config.yaml
EXPOSE 8080

CMD ["agentflow-server"]
