FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o game-server main.go
EXPOSE 8080
CMD ["./game-server","-env=prod"]
