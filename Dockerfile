FROM golang:1.17 as builder

WORKDIR /app

COPY . .

RUN go mod init example
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]