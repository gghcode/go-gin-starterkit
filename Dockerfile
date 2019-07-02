FROM golang:1.12.4-alpine AS builder
RUN apk add --no-cache git

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=builder /app/go-gin-starterkit /app/go-gin-starterkit

EXPOSE 8080
ENTRYPOINT [ "/app/go-gin-starterkit" ]
