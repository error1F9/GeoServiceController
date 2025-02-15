FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk update

WORKDIR /app

COPY --from=builder /app/cmd/main .

COPY ./hugo ./hugo

EXPOSE 8080

CMD ["/app/main"]