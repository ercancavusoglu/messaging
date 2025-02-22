FROM golang:1.21.1-alpine AS builder

WORKDIR /app

RUN apk update && \
    apk add --no-cache git && \
    rm -rf /var/cache/apk/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/app/main.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]