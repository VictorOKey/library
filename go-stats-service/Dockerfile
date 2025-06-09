FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8081
ENTRYPOINT ["./main"]
