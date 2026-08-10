# syntax=docker/dockerfile:1

FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/my-api \
    ./cmd/api


FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /app/bin/my-api /app/my-api

EXPOSE 8080

ENTRYPOINT ["/app/my-api"]