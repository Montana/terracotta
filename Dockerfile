FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o terracotta .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -s /bin/sh terracotta
WORKDIR /home/terracotta
COPY --from=builder /app/terracotta .
RUN chown terracotta:terracotta terracotta
USER terracotta
EXPOSE 8080 9090
ENTRYPOINT ["./terracotta"]
CMD ["-help"]

