FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o rpa-backend cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash mysql-client
WORKDIR /root/
COPY --from=builder /app/rpa-backend .
EXPOSE 8080
CMD ["./rpa-backend"]
