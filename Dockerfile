FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/server"]
