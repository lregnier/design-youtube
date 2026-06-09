FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./cmd/worker

FROM alpine:3.20
RUN apk add --no-cache ffmpeg
COPY --from=builder /bin/worker /worker
ENTRYPOINT ["/worker"]
