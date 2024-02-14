FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -o main .


FROM alpine

WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /app/main .

CMD ["./main"]
