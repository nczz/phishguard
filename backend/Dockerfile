FROM golang:1.24-alpine AS builder
ARG TARGET=api
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/phishguard ./cmd/${TARGET}

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/phishguard /usr/local/bin/phishguard
ENTRYPOINT ["phishguard"]
