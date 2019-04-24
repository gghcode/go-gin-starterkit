FROM golang:1.12.4-alpine AS builder
ENV GO111MODULE=on

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN apk add --no-cache git && go mod download && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=builder /app/go-gin-starterkit /app/go-gin-starterkit

EXPOSE 8080
ENTRYPOINT [ "/app/go-gin-starterkit" ]
