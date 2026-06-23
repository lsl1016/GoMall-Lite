FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gomall-api ./cmd/main.go

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/gomall-api /app/gomall-api
EXPOSE 8080
CMD ["/app/gomall-api"]
