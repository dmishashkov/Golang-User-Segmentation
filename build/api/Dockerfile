FROM golang:latest AS builder

WORKDIR /app

COPY . ./

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /app/server



FROM alpine:3

COPY --from=builder /app/server /bin/server

ENTRYPOINT ["/bin/server"]