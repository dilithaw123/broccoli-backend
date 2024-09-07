FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go build -o ./server ./cmd/web/main.go
CMD ["/app/server"]

