ARG GO_IMAGE=golang:1.23-bookworm
ARG ALPINE_IMAGE=debian:bookworm-slim

FROM ${GO_IMAGE} AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gomall-api ./cmd/main.go

FROM ${ALPINE_IMAGE}

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/gomall-api /app/gomall-api
EXPOSE 8080
CMD ["/app/gomall-api"]
