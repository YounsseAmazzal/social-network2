
FROM golang:1.27-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main ./main
COPY config.json ./config.json
COPY sqlinteraction/migrations ./sqlinteraction/migrations
COPY static ./static

RUN mkdir -p Db

EXPOSE 8081

CMD ["./main"]
